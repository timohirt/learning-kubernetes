package cmd

import (
	"fmt"
	"kthw/cmd/common"
	"kthw/cmd/infra/server"
	"kthw/cmd/infra/sshkey"
	"strings"

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
		common.WhenErrPrintAndExit(err)
		err = viper.WriteConfig()
		common.WhenErrPrintAndExit(err)

		fmt.Printf("SSH key '%s' successfully added to config.\n", name)
	}}

var addServerCommand = &cobra.Command{
	Use:   "add-server <name> <roles>",
	Short: "Adds a new server to the config file.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("Expected exactly two arguments, but found '%d'", len(args))
		}

		for _, role := range strings.Split(args[1], ",") {
			if server.IsValidRole(role) != nil {
				validRoles := strings.Join(server.AllValidRoles(), ", ")
				return fmt.Errorf("'%s' is not a valid role. Valid roles are: %s", role, validRoles)
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]
		roles := strings.Split(args[1], ",")

		err := server.AddServer(serverName, roles)
		common.WhenErrPrintAndExit(err)

		err = viper.WriteConfig()
		common.WhenErrPrintAndExit(err)

		fmt.Printf("Server %s successfully added to config.\n", serverName)
	}}

func configCommands() *cobra.Command {
	configCommand.AddCommand(initConfCommand)
	configCommand.AddCommand(addServerCommand)
	configCommand.AddCommand(addSSHKeyCommand)
	return configCommand
}
