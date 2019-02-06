package cmd

import (
	"fmt"
	"kthw/cmd/common"
	"kthw/cmd/server"
	"kthw/cmd/sshkey"
	"log"

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
		server.SetHCloudServerDefaults()
		err := viper.WriteConfig()
		common.WhenErrPrintAndExit(err)
	}}

var addSSHKeyCommand = &cobra.Command{
	Use:   "add-ssh-key <name> <file>",
	Short: "Adds a SSH public key to the config which will be set as authorized key of root user of created servers",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		file := args[1]
		err := sshkey.AddSSHPublicKeyToConfig(name, file)
		if err != nil {
			log.Panicf("Error while adding SSH key: %s", err)
		}

		err = viper.WriteConfig()
		if err != nil {
			log.Panicf("Error while writing SSH key to config file: %s", err)
		}
		fmt.Printf("SSH key '%s' successfully added to config.\n", name)
	}}

var addServerCommand = &cobra.Command{
	Use:   "add-server <name>",
	Short: "Adds a new server to the config file.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]
		server.AddServer(serverName)
		err := viper.WriteConfig()
		common.WhenErrPrintAndExit(err)
	}}

func configCommands() *cobra.Command {
	configCommand.AddCommand(initConfCommand)
	configCommand.AddCommand(addServerCommand)
	configCommand.AddCommand(addSSHKeyCommand)
	return configCommand
}
