package cmd

import (
	"fmt"
	"kthw/certs"
	"kthw/cmd/common"
	"os"

	"github.com/spf13/cobra"
)

var certsCommand = &cobra.Command{Use: "certs", Short: "Provision a certificate authority and generate certificates"}

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

var genEtcdClientCertificateCommand = &cobra.Command{
	Use:   "gen-etcd-client-cert",
	Short: "Generates etcd client certificate",
	Run: func(cmd *cobra.Command, args []string) {
		certGenerator := newCertificateGenerator()
		generateAndWriteEtcdClientCert(certGenerator)
	}}

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
	return certsCommand
}
