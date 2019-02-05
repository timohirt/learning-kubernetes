package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

// Execute runs commands child commands
func Execute() {
	viper.SetConfigFile(defaultConfigFile)
	viper.ReadInConfig()
	rootCmd.AddCommand(configCommands(), certsCommands(), provisionCommands())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
