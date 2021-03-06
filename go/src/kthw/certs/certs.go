package certs

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"path"
	"time"

	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
)

// WriteCert defines methods implemented by all cert to writeit to disk
type WriteCert interface {
	Write() error
}

// CertPaths defines accessors to get public and private key file paths
type CertPaths interface {
	PrivateKeyPath() string
	PublicKeyPath() string
}

// Cert has public and private keys, and base directory where both are stored
type cert struct {
	BaseDir         string
	PrivateKeyBytes []byte
	PublicKeyBytes  []byte
	CertPaths
	WriteCert
}

func writeCert(certPaths CertPaths, baseDir string, privateKeyBytes []byte, publicKeyBytes []byte) error {
	err := ensureDirectoryExists(baseDir)
	if err != nil {
		return err
	}
	err = writeToFile(privateKeyBytes, certPaths.PrivateKeyPath())
	if err != nil {
		return err
	}
	return writeToFile(publicKeyBytes, certPaths.PublicKeyPath())
}

// EtcdCert represents private and public key of etcd cert. This
// certificate is used for SSL.
type EtcdCert cert

// PrivateKeyPath gets the path to the admin private key file
func (e *EtcdCert) PrivateKeyPath() string { return path.Join(e.BaseDir, etcdKeyFileName) }

// PublicKeyPath gets the path to the admin private key file
func (e *EtcdCert) PublicKeyPath() string { return path.Join(e.BaseDir, etcdCertFileName) }

// EtcdClientCert private and public key used by clients to access
// a etcd cluster which uses client auth.
type EtcdClientCert cert

// PrivateKeyPath gets the path to the admin private key file
func (e *EtcdClientCert) PrivateKeyPath() string { return path.Join(e.BaseDir, etcdClientKeyFileName) }

// PublicKeyPath gets the path to the admin private key file
func (e *EtcdClientCert) PublicKeyPath() string { return path.Join(e.BaseDir, etcdClientCertFileName) }

// Write writes public and private key of a cert to files
func (e *EtcdClientCert) Write() error {
	return writeCert(e, e.BaseDir, e.PrivateKeyBytes, e.PublicKeyBytes)
}

// CertGenerator generates certificates using a CA
type CertGenerator struct {
	CA          *CA
	caKey       crypto.Signer
	caCert      *x509.Certificate
	signingConf *config.Signing
	certsConf   Config
	GeneratesCerts
}

// GeneratesCerts gets a CA and generates different certificates
type GeneratesCerts interface {
	GetCA() *CA
	GenEtcdCertificate(hosts []string) (*EtcdCert, error)
	GenEtcdClientCertificate() (*EtcdClientCert, error)
}

// NewCertGenerator creates a CertGenerator using given CACerts
func NewCertGenerator(ca *CA, certsConf Config) (*CertGenerator, error) {
	if ca == nil {
		return nil, fmt.Errorf("CACerts not porpery initiated. Either InitCA or LoadCA")
	}

	caKey, err := helpers.ParsePrivateKeyPEM(ca.KeyBytes)
	if err != nil {
		return nil, fmt.Errorf("Error while CA parsing private key: %s", err)
	}

	caCert, err := helpers.ParseCertificatePEM(ca.CertBytes)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing CA certificate: %s", err)
	}

	expiryDuration, _ := time.ParseDuration(signingExpiry)
	signingConf := &config.Signing{
		Profiles: map[string]*config.SigningProfile{
			signingProfile: {
				Usage:  signingUsages,
				Expiry: expiryDuration,
			},
		},
		Default: &config.SigningProfile{
			Expiry: expiryDuration,
		},
	}
	return &CertGenerator{
		CA:          ca,
		caKey:       caKey,
		caCert:      caCert,
		signingConf: signingConf,
		certsConf:   certsConf}, nil
}

// LoadCertGenerator loads a existing CA and creates a CertGenerator.
func LoadCertGenerator() (*CertGenerator, error) {
	ca, err := NewDefaultCertificateLoader().LoadCA()
	if err != nil {
		return nil, fmt.Errorf("Error while loading CA. %s", err)
	}

	conf := ReadConfig()
	certGenerator, err := NewCertGenerator(ca, conf)
	if err != nil {
		return nil, fmt.Errorf("Error while creating certificate generator: %s", err)
	}
	return certGenerator, nil
}

// GetCA returns the CA used to generate certificates
func (c *CertGenerator) GetCA() *CA { return c.CA }

const noHostname string = ""

// GenEtcdCertificate generates a etcd server certificate using the CA og CertGenerator.
func (c *CertGenerator) GenEtcdCertificate(hosts []string) (*EtcdCert, error) {
	req := &csr.CertificateRequest{
		CN:         etcdCN,
		KeyRequest: &csr.BasicKeyRequest{A: keyAlgo, S: keySize},
		Names:      []csr.Name{certName(etcdO)},
		Hosts:      hosts}
	privateKeyBytes, publicKeyBytes, _ := c.genPrivateAndPublicKey(req, noHostname)
	etcdCert := &EtcdCert{BaseDir: c.certsConf.BaseDir, PrivateKeyBytes: privateKeyBytes, PublicKeyBytes: publicKeyBytes}
	return etcdCert, nil
}

// GenEtcdClientCertificate generates a etcd client certificate using the CA og CertGenerator.
func (c *CertGenerator) GenEtcdClientCertificate() (*EtcdClientCert, error) {
	req := &csr.CertificateRequest{
		CN:         noHostname,
		KeyRequest: &csr.BasicKeyRequest{A: keyAlgo, S: keySize},
		Names:      []csr.Name{certName(etcdO)}}
	privateKeyBytes, publicKeyBytes, _ := c.genPrivateAndPublicKey(req, noHostname)
	etcdCert := &EtcdClientCert{BaseDir: c.certsConf.BaseDir, PrivateKeyBytes: privateKeyBytes, PublicKeyBytes: publicKeyBytes}
	return etcdCert, nil
}

func (c *CertGenerator) genPrivateAndPublicKey(req *csr.CertificateRequest, hostname string) (privateKeyBytes []byte, publicKeyBytes []byte, err error) {
	csrBytes, privateKeyBytes, err := c.genPrivateKey(req)
	if err != nil {
		return nil, nil, fmt.Errorf("Error while generating private key: %s", err)
	}
	publicKeyBytes, err = c.genPublicKey(csrBytes, privateKeyBytes, hostname)
	if err != nil {
		return nil, nil, fmt.Errorf("Error while generating public key: %s", err)
	}
	return privateKeyBytes, publicKeyBytes, nil
}

func (c *CertGenerator) genPrivateKey(req *csr.CertificateRequest) (csrBytes []byte, key []byte, err error) {
	gen := &csr.Generator{Validator: genkey.Validator}
	csrBytes, privateKeyBytes, err := gen.ProcessRequest(req)
	if err != nil {
		return nil, nil, fmt.Errorf("Error while creating private key %s", err)
	}
	return csrBytes, privateKeyBytes, nil
}

func (c *CertGenerator) genPublicKey(csrBytes []byte, privateKeyBytes []byte, hostname string) (cert []byte, err error) {
	caSigner, err := local.NewSigner(c.caKey, c.caCert, signer.DefaultSigAlgo(c.caKey), c.signingConf)
	if err != nil {
		return nil, fmt.Errorf("Error while creating CA signer: %s", err)
	}

	signReq := signer.SignRequest{
		Request: string(csrBytes),
		Hosts:   signer.SplitHosts(hostname),
		Profile: signingProfile,
	}
	return caSigner.Sign(signReq)
}
