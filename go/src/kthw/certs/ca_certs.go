package certs

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
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
func (c *CACerts) CNPrivateKeyFile() string { return path.Join(c.CABaseDir, caKeyFileName) }

// CNPublicKeyFile returns the path to CA private key PEM file.
func (c *CACerts) CNPublicKeyFile() string { return path.Join(c.CABaseDir, caCertFileName) }

// DefaultCACerts initializes Certs with default parameters.
func DefaultCACerts() *CACerts {
	caConf := &csr.CAConfig{PathLength: 0, PathLenZero: true, Expiry: signingExpiry}
	return &CACerts{
		CAKeySize: keySize,
		CABaseDir: "ca",
		cAConf:    caConf,
		cACsr: &csr.CertificateRequest{
			CN:         caCN,
			CA:         caConf,
			KeyRequest: &csr.BasicKeyRequest{A: keyAlgo, S: keySize},
			Names:      []csr.Name{certName(caCN)}}}
}

func (c *CACerts) generateCA() error {
	certBytes, _, keyBytes, err := initca.New(c.cACsr)

	if err != nil {
		return fmt.Errorf("Error generating CA certs: %s", err)
	}
	c.CA = &CA{CertBytes: certBytes, KeyBytes: keyBytes}
	return nil
}

// InitCa generates the CA public and private key and stores both in PEM
// format in directory 'ca' relative to working directory.
func (c *CACerts) InitCa() error {
	c.generateCA()
	err := ensureDirectoryExists(c.CABaseDir)
	if err != nil {
		return fmt.Errorf("Error while ensuring CA directories: %s", err)
	}

	err = writeToFile(c.CA.KeyBytes, c.CNPrivateKeyFile())
	if err != nil {
		return fmt.Errorf("Writing CA private key to file failed: %s", err)
	}

	err = writeToFile(c.CA.CertBytes, c.CNPublicKeyFile())
	if err != nil {
		return fmt.Errorf("Writing CA private key to file failed: %s", err)
	}

	return nil
}

// LoadCA loads private and public keys of CA from files.
func (c *CACerts) LoadCA() error {
	keyBytes, err := ioutil.ReadFile(c.CNPrivateKeyFile())
	if err != nil {
		return fmt.Errorf("Error reading private key file '%s': %s", c.CNPrivateKeyFile(), err)
	}
	certBytes, err := ioutil.ReadFile(c.CNPublicKeyFile())
	if err != nil {
		return fmt.Errorf("Error reading public key file '%s': %s", c.CNPublicKeyFile(), err)
	}
	c.CA = &CA{KeyBytes: keyBytes, CertBytes: certBytes}

	return nil
}
