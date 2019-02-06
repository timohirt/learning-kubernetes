package cmd

import (
	"kthw/cmd/common"
	"kthw/cmd/config"

	"github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

const (
	defaultConfigFile = "project.yaml"

	// ConfProjectNameKey is the config key of project name
	ConfProjectNameKey = "project.name"
)

var configCommand = &cobra.Command{Use: "config", Short: "Manage configuration"}

var initConfCommand = &cobra.Command{
	Use:   "new <project>",
	Short: "Creates a config file of a new project.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		viper.SetConfigFile(defaultConfigFile)
		projectName := args[0]
		viper.Set(ConfProjectNameKey, projectName)
		config.SetHCloudServerDefaults()
		err := viper.WriteConfig()
		common.WhenErrPrintAndExit(err)
	}}

func configCommands() *cobra.Command {
	configCommand.AddCommand(initConfCommand)
	configCommand.AddCommand(config.AddServerCommand)
	configCommand.AddCommand(addSSHKeyCommand)
	return configCommand
}
