package cmd

import (
	"fmt"
	"kthw/certs"
	"kthw/cmd/common"
	"kthw/cmd/hcloudclient"
	"kthw/cmd/infra/server"
	"kthw/cmd/infra/sshkey"
	"os"
	"strings"

	"github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

const (
	defaultConfigFile = "project.yaml"

	// ConfProjectNameKey is the config key of project name
	ConfProjectNameKey = "project.name"
)

var projectCommand = &cobra.Command{Use: "project", Short: "Create and manage configuration of a project"}

var newProjectCommand = &cobra.Command{
	Use:   "new <name> <ssh-public-key-file>",
	Short: "Creates a config file of a new project and sets up everything required to provision K8s clusters.",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if APIToken == "" {
			fmt.Println("ApiToken not found. Make sure you set the --apiToken flag")
			os.Exit(1)
		}

		viper.SetConfigFile(defaultConfigFile)
		projectName := args[0]
		sshPublicKeyFilePath := args[1]
		viper.Set(ConfProjectNameKey, projectName)
		server.SetHCloudServerDefaults()
		fmt.Println("Initialised project.yaml with defaults.")

		sshPublicKey, err := sshkey.AddSSHPublicKeyToConfig(projectName, sshPublicKeyFilePath)
		common.WhenErrPrintAndExit(err)
		fmt.Println("Added SSH key to config.")

		hcloudClient := hcloudclient.NewHCloudClient(APIToken)
		updatedConfig := sshkey.CreateSSHKey(*sshPublicKey, hcloudClient)
		updatedConfig.WriteToConfig()
		fmt.Println("SSH key created in hcloud.")

		certsConf := certs.InitDefaultConfig()
		caCerts := certs.DefaultCACerts(certsConf.BaseDir)
		err = caCerts.InitCa()
		if err != nil {
			fmt.Printf("Error while initiation CA: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("CA private and public keys generated and stored in %s", caCerts.CABaseDir)
		fmt.Println("Initialised PKI infrastructure.")

		err = viper.WriteConfig()
		common.WhenErrPrintAndExit(err)
	}}

var addServerCommand = &cobra.Command{
	Use:   "add-server <name> <roles>",
	Short: "Adds a new server to the config file.",
	Long:  "Pick a random name for this server. Valid roles are controller, worker and etcd.",
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

func projectCommands() *cobra.Command {
	projectCommand.AddCommand(newProjectCommand)
	projectCommand.AddCommand(addServerCommand)
	return projectCommand
}
