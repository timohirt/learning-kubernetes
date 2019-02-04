package certs

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"time"

	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/helpers"
	"github.com/cloudflare/cfssl/signer"
	"github.com/cloudflare/cfssl/signer/local"
)

// Cert has public and private keys, and base directory where both are stored
type cert struct {
	BaseDir         string
	PrivateKeyBytes []byte
	PublicKeyBytes  []byte
}

// AdminClientCert has admin client public and private
type AdminClientCert cert

// CertGenerator generates certificates using a CA
type CertGenerator struct {
	caCerts     *CACerts
	caKey       crypto.Signer
	caCert      *x509.Certificate
	signingConf *config.Signing
}

// NewCertGenerator creates a CertGenerator using given CACerts
func NewCertGenerator(ca *CACerts) (*CertGenerator, error) {
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
		caCerts:     ca,
		caKey:       caKey,
		caCert:      caCert,
		signingConf: signingConf}, nil
}

const noHostname string = ""

// GenAdminClientCertificate generates a admin client certificate using the CA og CertGenerator.
func (c *CertGenerator) GenAdminClientCertificate() (*AdminClientCert, error) {
	req := &csr.CertificateRequest{
		CN:         adminClientCN,
		KeyRequest: &csr.BasicKeyRequest{A: keyAlgo, S: keySize},
		Names:      []csr.Name{certName(adminClientO)}}
	privateKeyBytes, publicKeyBytes, _ := c.genPrivateAndPublicKey(req, noHostname)
	adminClientCert := &AdminClientCert{BaseDir: certsBaseDir, PrivateKeyBytes: privateKeyBytes, PublicKeyBytes: publicKeyBytes}
	return adminClientCert, nil
}

func (c *CertGenerator) genPrivateAndPublicKey(req *csr.CertificateRequest, hostname string) (privateKeyBytes []byte, publicKeyBytes []byte, err error) {
	csrBytes, privateKeyBytes, err := c.genPrivateKey(req)
	if err != nil {
		return nil, nil, fmt.Errorf("Error while generating private key: %s", err)
	}
	publicKeyBytes, err = c.genPublicKey(csrBytes, privateKeyBytes, noHostname)
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
