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

// AdminClientCert has admin client public and private
type AdminClientCert cert

// PrivateKeyPath gets the path to the admin private key file
func (a *AdminClientCert) PrivateKeyPath() string { return path.Join(a.BaseDir, adminClientKeyFileName) }

// PublicKeyPath gets the path to the admin private key file
func (a *AdminClientCert) PublicKeyPath() string { return path.Join(a.BaseDir, adminClientCertFileName) }

// Write writes public and private key of a cert to files
func (a *AdminClientCert) Write() error {
	return writeCert(a, a.BaseDir, a.PrivateKeyBytes, a.PublicKeyBytes)
}

// LoadAdminClientCert loads private and public key of a admin client certificate
func LoadAdminClientCert(config Config) (*AdminClientCert, error) {
	cert := AdminClientCert{
		BaseDir: config.BaseDir}

	privateKeyBytes, err := readFromFile(cert.PrivateKeyPath())
	if err != nil {
		return nil, fmt.Errorf("Could not load private key: '%s'", err)
	}
	cert.PrivateKeyBytes = privateKeyBytes

	publicKeyBytes, err := readFromFile(cert.PublicKeyPath())
	if err != nil {
		return nil, fmt.Errorf("Could not load private key: '%s'", err)
	}
	cert.PublicKeyBytes = publicKeyBytes

	return &cert, nil
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

// CertGenerator generates certificates using a CA
type CertGenerator struct {
	CACerts     *CACerts
	caKey       crypto.Signer
	caCert      *x509.Certificate
	signingConf *config.Signing
	certsConf   Config
	GeneratesCerts
}

// GeneratesCerts gets a CA and generates different certificates
type GeneratesCerts interface {
	GetCA() *CA
	GenAdminClientCertificate() (*AdminClientCert, error)
	GenEtcdCertificate(hosts []string) (*EtcdCert, error)
}

// NewCertGenerator creates a CertGenerator using given CACerts
func NewCertGenerator(ca *CACerts, certsConf Config) (*CertGenerator, error) {
	if ca.CA == nil {
		return nil, fmt.Errorf("CACerts not porpery initiated. Either InitCA or LoadCA")
	}

	caKey, err := helpers.ParsePrivateKeyPEM(ca.CA.KeyBytes)
	if err != nil {
		return nil, fmt.Errorf("Error while CA parsing private key: %s", err)
	}

	caCert, err := helpers.ParseCertificatePEM(ca.CA.CertBytes)
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
		CACerts:     ca,
		caKey:       caKey,
		caCert:      caCert,
		signingConf: signingConf,
		certsConf:   certsConf}, nil
}

// LoadCertGenerator loads a existing CA and creates a CertGenerator.
func LoadCertGenerator() (*CertGenerator, error) {
	conf := ReadConfig()
	caCerts, err := LoadCACerts(conf)
	if err != nil {
		return nil, fmt.Errorf("Error while loading CA. %s", err)
	}

	certGenerator, err := NewCertGenerator(caCerts, conf)
	if err != nil {
		return nil, fmt.Errorf("Error while creating certificate generator: %s", err)
	}
	return certGenerator, nil
}

// GetCA returns the CA used to generate certificates
func (c *CertGenerator) GetCA() *CA { return c.CACerts.CA }

const noHostname string = ""

// GenAdminClientCertificate generates a admin client certificate using the CA og CertGenerator.
func (c *CertGenerator) GenAdminClientCertificate() (*AdminClientCert, error) {
	req := &csr.CertificateRequest{
		CN:         adminClientCN,
		KeyRequest: &csr.BasicKeyRequest{A: keyAlgo, S: keySize},
		Names:      []csr.Name{certName(adminClientO)}}
	privateKeyBytes, publicKeyBytes, _ := c.genPrivateAndPublicKey(req, noHostname)
	adminClientCert := &AdminClientCert{BaseDir: c.certsConf.BaseDir, PrivateKeyBytes: privateKeyBytes, PublicKeyBytes: publicKeyBytes}
	return adminClientCert, nil
}

// GenEtcdCertificate generates a admin client certificate using the CA og CertGenerator.
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
