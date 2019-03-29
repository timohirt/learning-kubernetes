package cmd

import (
	"fmt"
	"kthw/certs"
	"kthw/cmd/common"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"os"

	"github.com/spf13/cobra"
)

var provisionCommand = &cobra.Command{
	Use:   "provision",
	Short: "Commands for provisioning servers and clusters"}

var createServerCommand = &cobra.Command{
	Use:   "server <name>",
	Short: "Creates a server previously added to the config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]
		serverConfig := server.FromConfig(serverName)
		createServerAndUpdateConfig(&serverConfig)

		fmt.Printf("Server %s successfully created.\n", serverName)
	}}

var configureWireguardCommand = &cobra.Command{
	Use:   "network",
	Short: "Generates wireguard config and establishes private overlay network",
	Run: func(cmd *cobra.Command, args []string) {
		sshClient := sshconnect.NewSSHConnect(Verbose)
		serverConfigs, err := server.AllFromConfig()
		if err != nil {
			fmt.Printf("Error while loading servers from configuration: %s\n", err)
			os.Exit(1)
		}

		setupWireguardAndUpdateConfig(serverConfigs, sshClient)
	}}

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
		common.WhenErrPrintAndExit(err)

		installEtcd(serverConfigs, sshClient, certGenerator)
	}}

var installKubernetesControllerCommand = &cobra.Command{
	Use:   "k8s-controller",
	Short: "Generate config, upload certificates and install controller on node",
	Run: func(cmd *cobra.Command, args []string) {
		sshClient := sshconnect.NewSSHConnect(Verbose)
		certLoader := certs.NewDefaultCertificateLoader()
		certGenerator, err := certs.LoadCertGenerator()
		common.WhenErrPrintAndExit(err)
		serverConfigs, err := server.AllFromConfig()
		if err != nil {
			fmt.Printf("Error while loading servers from configuration: %s\n", err)
			os.Exit(1)
		}

		installKubernetesController(serverConfigs, sshClient, certLoader, certGenerator)
	}}

func provisionCommands() *cobra.Command {
	provisionCommand.AddCommand(createServerCommand)
	provisionCommand.AddCommand(configureWireguardCommand)
	provisionCommand.AddCommand(installEtcdCommand)
	provisionCommand.AddCommand(installKubernetesControllerCommand)
	provisionCommand.AddCommand(certsCommands())
	return provisionCommand
}
