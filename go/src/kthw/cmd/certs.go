package cmd

import (
	"kthw/certs"
	"log"

	"github.com/spf13/cobra"
)

var certsCommand = &cobra.Command{Use: "certs", Short: "Create CA, Server and Client certificates"}

var initCACommand = &cobra.Command{
	Use:   "init-ca",
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

var genAdminCertificateCommand = &cobra.Command{
	Use:   "gen-admin-cert",
	Short: "Generates admin certificate",
	Run: func(cmd *cobra.Command, args []string) {
		caCerts, err := certs.LoadCACerts()
		if err != nil {
			log.Fatal(err)
		}
		certGenerator, err := certs.NewCertGenerator(caCerts)
		if err != nil {
			log.Fatalf("Error while generating admin cert: %s", err)
		}
		adminCert, err := certGenerator.GenAdminClientCertificate()
		if err != nil {
			log.Fatalf("Error while generating admin cert: %s", err)
		}
		adminCert.WriteCert()
	}}

func certsCommands() *cobra.Command {
	certsCommand.AddCommand(initCACommand)
	certsCommand.AddCommand(genAdminCertificateCommand)
	return certsCommand
}
