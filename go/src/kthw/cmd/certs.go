package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
)

/*type CertRequests struct {
	CAConfig *csr.CAConfig
	Request  *csr.CertificateRequest
}

func (cs CertRequests) New() CertRequests {

}*/

type Certs struct {
	CAKeySize int
	CABaseDir string
	CAConf    *csr.CAConfig
	CACsr     *csr.CertificateRequest
	CA        *CA
}

// DefaultCerts initializes Certs with default parameters.
func DefaultCerts() *Certs {
	caConf := &csr.CAConfig{PathLength: 0, PathLenZero: true, Expiry: "8760h"}
	keySize := 2048
	return &Certs{
		CAKeySize: keySize,
		CABaseDir: "ca",
		CAConf:    caConf,
		CACsr: &csr.CertificateRequest{
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

// CA holds the certificates used to generate SSL certificates.
type CA struct {
	CertBytes []byte
	KeyBytes  []byte
}

func (c *Certs) generateCA() error {
	certBytes, _, keyBytes, err := initca.New(c.CACsr)

	if err != nil {
		return fmt.Errorf("Error generating CA certs: %s", err)
	}
	c.CA = &CA{CertBytes: certBytes, KeyBytes: keyBytes}
	return nil
}

func (c *Certs) ensureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Certs) writeToFileOrDie(cert []byte, file string) {
	err := ioutil.WriteFile(file, cert, 0644)
	if err != nil {
		panic(err)
	}
}

func (c *Certs) InitCa() error {
	c.generateCA()
	err := c.ensureDirectoryExists(c.CABaseDir)
	if err != nil {
		return fmt.Errorf("Error while ensuring CA directories: %s", err)
	}

	c.writeToFileOrDie(c.CA.KeyBytes, path.Join(c.CABaseDir, "ca-key.pem"))

	return nil
}
