package cmd

import (
	"kthw/certs"
	"kthw/cmd/common"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"sync"

	"github.com/spf13/cobra"
)

var installCommand = &cobra.Command{
	Use:   "install",
	Short: "Install infrastructure and kubernetes cluster with one command"}

var kubernetesCluster = &cobra.Command{
	Use:   "k8s-non-ha",
	Short: "Install a non HA cluster",
	Run: func(cmd *cobra.Command, args []string) {
		sshClient := sshconnect.NewSSHConnect(Verbose)
		certGenerator, err := certs.LoadCertGenerator()
		common.WhenErrPrintAndExit(err)
		certLoader := certs.NewDefaultCertificateLoader()

		serverConfigs, _ := server.AllFromConfig()
		for _, conf := range serverConfigs {
			createServerAndUpdateConfig(conf)
		}

		var waitGroup sync.WaitGroup

		for _, conf := range serverConfigs {
			waitGroup.Add(1)
			go waitForCloudInitCompleted(&waitGroup, conf, sshClient)
		}

		waitGroup.Wait()

		setupWireguardAndUpdateConfig(serverConfigs, sshClient)
		installEtcd(serverConfigs, sshClient, certGenerator)
		installKubernetesController(serverConfigs, sshClient, certLoader, certGenerator)
	}}

func installCommands() *cobra.Command {
	installCommand.PersistentFlags().StringVarP(&APIToken, "apiToken", "a", "", "API token for access to hcloud (required)")
	installCommand.AddCommand(kubernetesCluster)
	return installCommand
}
