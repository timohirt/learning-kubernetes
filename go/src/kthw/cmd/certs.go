package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	"github.com/spf13/cobra"
)

// CACerts stores CA configuration and manages certificates.
type CACerts struct {
	CAKeySize int
	CABaseDir string
	cACsr     *csr.CertificateRequest
	cAConf    *csr.CAConfig
	CA        *CA
}

// CA holds the certificates used to generate SSL certificates.
type CA struct {
	CertBytes []byte
	KeyBytes  []byte
}

// CNPrivateKeyFile returns the path to CA private key PEM file.
func (c *CACerts) CNPrivateKeyFile() string { return path.Join(c.CABaseDir, "ca-key.pem") }

// CNPublicKeyFile returns the path to CA private key PEM file.
func (c *CACerts) CNPublicKeyFile() string { return path.Join(c.CABaseDir, "ca.pem") }

// DefaultCACerts initializes Certs with default parameters.
func DefaultCACerts() *CACerts {
	caConf := &csr.CAConfig{PathLength: 0, PathLenZero: true, Expiry: "8760h"}
	keySize := 2048
	return &CACerts{
		CAKeySize: keySize,
		CABaseDir: "ca",
		cAConf:    caConf,
		cACsr: &csr.CertificateRequest{
			CN:         "Kubernetes",
			CA:         caConf,
			KeyRequest: &csr.BasicKeyRequest{A: "rsa", S: keySize},
			Names:      []csr.Name{certName("Kubernetes")}}}
}

func certName(o string) csr.Name {
	return csr.Name{
		C:  "DE",
		L:  "Mainz",
		O:  o,
		OU: "Learning Kubernetes",
		ST: "RLP"}
}

func (c *CACerts) generateCA() error {
	certBytes, _, keyBytes, err := initca.New(c.cACsr)

	if err != nil {
		return fmt.Errorf("Error generating CA certs: %s", err)
	}
	c.CA = &CA{CertBytes: certBytes, KeyBytes: keyBytes}
	return nil
}

func (c *CACerts) ensureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *CACerts) writeToFileOrDie(cert []byte, file string) {
	err := ioutil.WriteFile(file, cert, 0644)
	if err != nil {
		panic(err)
	}
}

// InitCa generates the CA public and private key and stores both in PEM
// format in directory 'ca' relative to working directory.
func (c *CACerts) InitCa() error {
	c.generateCA()
	err := c.ensureDirectoryExists(c.CABaseDir)
	if err != nil {
		return fmt.Errorf("Error while ensuring CA directories: %s", err)
	}

	c.writeToFileOrDie(c.CA.KeyBytes, c.CNPrivateKeyFile())
	c.writeToFileOrDie(c.CA.CertBytes, c.CNPublicKeyFile())

	return nil
}

var certsCommand = &cobra.Command{Use: "certs"}

var initCACommand = &cobra.Command{Use: "init-ca",
	Short: "Generates CA public and private key",
	Run: func(cmd *cobra.Command, args []string) {
		caCerts := DefaultCACerts()
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
