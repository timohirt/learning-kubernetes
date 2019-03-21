package cmd

import (
	"fmt"
	"kthw/certs"
	"kthw/cmd/cluster/etcd"
	"kthw/cmd/cluster/kube"
	"kthw/cmd/common"
	"kthw/cmd/hcloudclient"
	"kthw/cmd/infra/network"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"sync"
	"time"

	"github.com/spf13/viper"
)

func createServerAndUpdateConfig(config *server.Config) {
	fmt.Printf("Creating server %s at Hetzner cloud\n", config.Name)
	hcloudClient := hcloudclient.NewHCloudClient(APIToken)
	err := server.Create(config, hcloudClient)
	common.WhenErrPrintAndExit(err)

	config.UpdateConfig()
	viper.WriteConfig()
}

func setupWireguardAndUpdateConfig(configs []*server.Config, sshClient sshconnect.SSHOperations) {
	fmt.Println("Setting up private overlay network")
	err := network.SetupWireguard(sshClient, configs)
	common.WhenErrPrintAndExit(err)

	for _, config := range configs {
		config.UpdateConfig()
		fmt.Printf("Wireguard set up on %s.\n", config.Name)
	}
	viper.WriteConfig()
}

func installEtcd(configs []*server.Config, sshClient sshconnect.SSHOperations, certGenerator certs.GeneratesCerts) {
	fmt.Println("Installing etcd")

	err := etcd.InstallOnHost(configs, sshClient, certGenerator)
	common.WhenErrPrintAndExit(err)
}

func installKubernetesController(configs []*server.Config, sshclient sshconnect.SSHOperations, certLoader certs.CertificateLoader) {
	fmt.Println("Installing kubernetes controller")

	err := kube.InstallOnHosts(configs, sshclient, certLoader)
	common.WhenErrPrintAndExit(err)
}

func waitForCloudInitCompleted(waitGroup *sync.WaitGroup, conf *server.Config, sshClient sshconnect.SSHOperations) {
	fmt.Printf("Waiting for %s to complete cloud-init\n", conf.Name)
	defer waitGroup.Done()

	cloudInitRunning := true
	for retries := 0; retries < 20 && cloudInitRunning; retries++ {
		if server.IsCloudInitCompleted(conf.PublicIP, sshClient) {
			fmt.Printf("%s completed cloud-init\n", conf.Name)
			cloudInitRunning = false
		} else {
			time.Sleep(time.Second * 20)
		}
	}
}
