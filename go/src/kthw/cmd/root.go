package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// APIToken used to authenticate with Hetzer Cloud API
var APIToken string

var rootCmd = &cobra.Command{
	Use:   "hk",
	Short: "hk helps to setup Kubernetes in Hetzner Cloud", Long: `See short`}

// Execute runs commands child commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&APIToken, "apiToken", "a", "", "Hetzner API token (required)")
	rootCmd.MarkFlagRequired("apiToken")
}
