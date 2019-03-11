package cmd

import (
	"fmt"
	"kthw/certs"
	"kthw/cmd/cluster/etcd"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"os"

	"github.com/spf13/cobra"
)

var clusterCommand = &cobra.Command{
	Use:   "install",
	Short: "Commands for installing K8S cluster"}

var installEtcdCommand = &cobra.Command{
	Use:   "etcd",
	Short: "Downloads and installs etcd",
	Run: func(cmd *cobra.Command, args []string) {
		sshClient := sshconnect.NewSSHConnect(Verbose)
		serverConfigs, err := server.AllFromConfig()
		if err != nil {
			fmt.Printf("Error while loading servers from configuration: %s\n", err)
			os.Exit(1)
		}

		certGenerator, err := certs.LoadCertGenerator()
		if err != nil {
			fmt.Printf("Error while creating certificate generator: %s\n", err)
			os.Exit(1)
		}

		err = etcd.InstallOnHost(serverConfigs, sshClient, certGenerator)
		if err != nil {
			fmt.Printf("Error while installing etcd: %s\n", err)
		}
	}}

func clusterCommands() *cobra.Command {
	clusterCommand.PersistentFlags().StringVarP(&APIToken, "apiToken", "a", "", "API token for access to hcloud (required)")
	clusterCommand.AddCommand(installEtcdCommand)
	return clusterCommand
}
