package certs

import (
	"fmt"
)

// CertificateLoader loads public and private keys
type CertificateLoader interface {
	LoadEtcdClientCert() (*EtcdClientCert, error)
	LoadCA() (*CA, error)
}

// DefaultCertificateLoader loads certificates from filesystem
type DefaultCertificateLoader struct {
	baseDir string
	CertificateLoader
}

// NewDefaultCertificateLoader creates a DefaultCertificateLoader using certificate base dir from certs config.
func NewDefaultCertificateLoader() *DefaultCertificateLoader {
	config := ReadConfig()
	return &DefaultCertificateLoader{baseDir: config.BaseDir}
}

// LoadEtcdClientCert loads etcd client certificate from filesystem.
func (d *DefaultCertificateLoader) LoadEtcdClientCert() (*EtcdClientCert, error) {
	cert := &EtcdClientCert{
		BaseDir: d.baseDir}

	privateKeyBytes, publicKeyBytes, err := d.loadPrivateAndPublicKey(cert.PrivateKeyPath(), cert.PublicKeyPath())
	if err != nil {
		return nil, err
	}

	cert.PrivateKeyBytes = privateKeyBytes
	cert.PublicKeyBytes = publicKeyBytes

	return cert, nil
}

// LoadCA loads CA certificate from filesystem.
func (d *DefaultCertificateLoader) LoadCA() (*CA, error) {
	caCert := DefaultCACerts(d.baseDir)
	err := caCert.LoadCA()
	if err != nil {
		return nil, err
	}

	return caCert.CA, nil
}

func (d *DefaultCertificateLoader) loadPrivateAndPublicKey(privateKeyPath string, publicKeyPath string) ([]byte, []byte, error) {
	privateKeyBytes, err := readFromFile(privateKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not load private key: '%s'", err)
	}

	publicKeyBytes, err := readFromFile(publicKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not load public key: '%s'", err)
	}

	return privateKeyBytes, publicKeyBytes, nil
}
