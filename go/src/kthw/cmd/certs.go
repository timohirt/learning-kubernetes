package cmd

import (
	"fmt"
	"kthw/certs"
	"kthw/cmd/common"
	"os"

	"github.com/spf13/cobra"
)

var certsCommand = &cobra.Command{Use: "certs", Short: "Create CA, Server and Client certificates"}

var initCACommand = &cobra.Command{
	Use:   "init-ca",
	Short: "Generates CA public and private key",
	Run: func(cmd *cobra.Command, args []string) {
		conf := certs.ReadConfig()
		caCerts := certs.DefaultCACerts(conf.BaseDir)
		err := caCerts.InitCa()
		if err != nil {
			fmt.Printf("Error while initiation CA: %s\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("CA private and public keys generated and stored in %s", caCerts.CABaseDir)
		}
	}}

var genAdminCertificateCommand = &cobra.Command{
	Use:   "gen-admin-cert",
	Short: "Generates admin certificate",
	Run: func(cmd *cobra.Command, args []string) {
		certGenerator := newCertificateGenerator()
		generateAndWriteAdminClientCert(certGenerator)
	}}

var genEtcdClientCertificateCommand = &cobra.Command{
	Use:   "gen-etcd-client-cert",
	Short: "Generates etcd client certificate",
	Run: func(cmd *cobra.Command, args []string) {
		certGenerator := newCertificateGenerator()
		generateAndWriteEtcdClientCert(certGenerator)
	}}

var genAllCertificatesCommand = &cobra.Command{
	Use:   "gen-all-certs",
	Short: "Generates all certificates for Kubernetes and etcd",
	Run: func(cmd *cobra.Command, args []string) {
		conf := certs.ReadConfig()
		fmt.Printf("Loading CA from %s.\n", conf.BaseDir)
		certLoader := certs.NewDefaultCertificateLoader()
		ca, err := certLoader.LoadCA()
		if err != nil {
			fmt.Printf("CA not found. Generating new CA.\n")
			caCerts := certs.DefaultCACerts(conf.BaseDir)
			err = caCerts.InitCa()
			if err != nil {
				fmt.Printf("Initializing CA failed. %s\n", err)
				os.Exit(1)
			}
			ca = caCerts.CA
		}

		certGenerator, err := certs.NewCertGenerator(ca, conf)
		if err != nil {
			fmt.Printf("Error while creating certificate generator: %s\n", err)
			os.Exit(1)
		}

		generateAndWriteAdminClientCert(certGenerator)
	}}

func generateAndWriteAdminClientCert(certGenerator *certs.CertGenerator) {
	fmt.Printf("Generating admin certificate.\n")
	adminCert, err := certGenerator.GenAdminClientCertificate()
	if err != nil {
		fmt.Printf("Error while generating admin cert: %s\n", err)
		os.Exit(1)
	}
	adminCert.Write()
}

func generateAndWriteEtcdClientCert(certGenerator *certs.CertGenerator) {
	fmt.Printf("Generating etcd client certificate.\n")
	etcdClientCert, err := certGenerator.GenEtcdClientCertificate()
	if err != nil {
		fmt.Printf("Error while generating etcd client cert: %s\n", err)
		os.Exit(1)
	}
	etcdClientCert.Write()
}

func newCertificateGenerator() *certs.CertGenerator {
	ca, err := certs.NewDefaultCertificateLoader().LoadCA()
	common.WhenErrPrintAndExit(err)

	conf := certs.ReadConfig()
	certGenerator, err := certs.NewCertGenerator(ca, conf)
	common.WhenErrPrintAndExit(err)
	return certGenerator
}

func certsCommands() *cobra.Command {
	certsCommand.AddCommand(initCACommand)
	certsCommand.AddCommand(genEtcdClientCertificateCommand)
	certsCommand.AddCommand(genAdminCertificateCommand)
	certsCommand.AddCommand(genAllCertificatesCommand)
	return certsCommand
}
