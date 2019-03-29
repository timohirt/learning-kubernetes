package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Long: "First create a project, add ssh-keys and some servers. Then use the install commands to install the cluster."}

// APIToken used to authenticate with Hetzer Cloud API.
var APIToken string

// Verbose controls the verbosity during command execution. If 'true' you'll see more output which might help for debugging.
var Verbose bool

// Execute runs commands child commands
func Execute() {
	viper.SetConfigFile(defaultConfigFile)
	viper.ReadInConfig()
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Set 'true' for more output.")
	rootCmd.PersistentFlags().StringVarP(&APIToken, "apiToken", "a", "", "API token for access to hcloud (required)")
	rootCmd.AddCommand(projectCommands(), provisionCommands(), installCommands())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
