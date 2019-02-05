package cmd

import (
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
		serverConfig := serverConfigFromConfig(serverName)
		hcloudClient := NewHCloudClient(APIToken)
		updatedConfig, err := CreateServer(serverConfig, hcloudClient)
		whenErrPrintAndExit(err)

		updatedConfig.UpdateConfig()
		viper.WriteConfig()
	}}

func provisionCommands() *cobra.Command {
	createServerCommand.Flags().StringVarP(&APIToken, "apiToken", "a", "", "API token for access to hcloud (required)")
	provisionCommand.AddCommand(createServerCommand)
	return provisionCommand
}
