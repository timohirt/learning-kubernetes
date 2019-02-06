package cmd

import (
	"fmt"
	"kthw/cmd/common"
	"kthw/cmd/config"
	"kthw/cmd/hcloudclient"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var provisionCommand = &cobra.Command{
	Use:   "provision",
	Short: "Commands for provisioning servers"}

var createServerCommand = &cobra.Command{
	Use:   "create <name>",
	Short: "Creates a server previously added to the config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]
		serverConfig := config.ServerConfigFromConfig(serverName)
		hcloudClient := hcloudclient.NewHCloudClient(APIToken)
		updatedConfig, err := createServer(serverConfig, hcloudClient)
		common.WhenErrPrintAndExit(err)

		updatedConfig.UpdateConfig()
		viper.WriteConfig()
	}}

var createSSHKeysCommand = &cobra.Command{
	Use:   "ssh-keys",
	Short: "Reads ssh key from config and creates in in hcloud",
	Run: func(cmd *cobra.Command, args []string) {
		key, err := common.ReadSSHPublicKeyFromConf()
		common.WhenErrPrintAndExit(err)
		if APIToken == "" {
			fmt.Println("ApiToken not found. Make sure you set the --apitoken flag")
			os.Exit(1)
		}
		hcloudClient := hcloudclient.NewHCloudClient(APIToken)
		updatedConfig := createSSHKey(*key, hcloudClient)

		updatedConfig.WriteToConfig()
		viper.WriteConfig()
		fmt.Println("SSH key created at hcloud.")
	}}

func provisionCommands() *cobra.Command {
	provisionCommand.PersistentFlags().StringVarP(&APIToken, "apiToken", "a", "", "API token for access to hcloud (required)")
	provisionCommand.AddCommand(createServerCommand)
	provisionCommand.AddCommand(createSSHKeysCommand)
	return provisionCommand
}
