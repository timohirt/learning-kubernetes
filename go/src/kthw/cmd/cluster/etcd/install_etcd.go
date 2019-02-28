package etcd

import (
	"fmt"
	"kthw/cmd/common"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"os"
)

func InstallOnHost(hostConfigs []server.Config, ssh sshconnect.SSHOperations) error {
	etcdHosts := selectEtcdHosts(hostConfigs)
	if len(etcdHosts) == 0 {
		return fmt.Errorf("List of provided hosts didn't contain a host with role etcd")
	}

	if len(etcdHosts) > 1 {
		return fmt.Errorf("Several hosts have role etcd. Currently, only single node etcd clusters are allowed")
	}

	for _, etcdHost := range etcdHosts {
		host := etcdHost.PublicIP
		commands := &sshconnect.Commands{
			Commands: []sshconnect.Command{
				downloadEtcd(host),
				unpackAndInstall(host),
				uploadSystemdService(etcdHost)},
			LogOutput: true}
		err := ssh.RunCmds(commands)
		if err != nil {
			return err
		}
	}
	return nil
}

func selectEtcdHosts(hostConfigs []server.Config) []server.Config {
	var etcdHosts []server.Config
	for _, host := range hostConfigs {
		if common.ArrayContains(host.Roles, "etcd") {
			etcdHosts = append(etcdHosts, host)
		}
	}
	return etcdHosts
}

func uploadSystemdService(hostConfig server.Config) *sshconnect.CopyFileCommand {
	params := EtcdSystemdServiceParameters{
		PrivateIP: hostConfig.PrivateIP,
		NodeName:  hostConfig.Name}
	systemdService, err := GenerateEtcdSystemdService(params)
	if err != nil {
		fmt.Printf("Error generating systemd service! %s\n", err)
		os.Exit(1)
	}

	return &sshconnect.CopyFileCommand{
		Host:        hostConfig.PublicIP,
		FileContent: systemdService,
		FilePath:    "/etc/systemd/system/etcd.service",
		Description: "Copy etcd systemd service to host"}
}

func downloadEtcd(host string) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: "curl -L https://github.com/etcd-io/etcd/releases/download/v3.3.12/etcd-v3.3.12-linux-amd64.tar.gz -o /tmp/etcd-v3.3.12.tar.gz",
		Host:        host,
		Description: "Download etcd binary"}
}

func unpackAndInstall(host string) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		CommandLine: "tar xzf /tmp/etcd-v3.3.12.tar.gz -C /tmp && mv /tmp/etcd-v3.3.12*/etcd* /usr/local/bin/",
		Host:        host,
		Description: "Untar etcd archive and copy to /usr/local/bin"}
}
