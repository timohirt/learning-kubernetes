package cmd

import (
	"kthw/certs"
	"log"

	"github.com/spf13/cobra"
)

var certsCommand = &cobra.Command{Use: "certs"}

var initCACommand = &cobra.Command{Use: "init-ca",
	Short: "Generates CA public and private key",
	Run: func(cmd *cobra.Command, args []string) {
		caCerts := certs.DefaultCACerts()
		err := caCerts.InitCa()
		if err != nil {
			log.Fatalf("Error while initiation CA: %s", err)
		} else {
			log.Printf("CA private and public keys generated and stored in %s", caCerts.CABaseDir)
		}
	}}

func init() {
	certsCommand.AddCommand(initCACommand)
	rootCmd.AddCommand(certsCommand)
}
