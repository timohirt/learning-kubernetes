package certs_test

import (
	"kthw/certs"
	"testing"
)

func TestCertGeneratorGenAdminClientCert(t *testing.T) {
	caCerts, _ := helperCreateDefaultCACerts(t)
	caCerts.InitCa()
	helperEnsureCaCertsInitialized(t, caCerts)
	certGenerator, err := certs.NewCertGenerator(caCerts)
	helperFailIfErr(t, "Error while creating CertGenerator", err)

	adminClientCert, err := certGenerator.GenAdminClientCertificate()
	helperFailIfErr(t, "Error creating admin client certificate", err)

	if adminClientCert.PrivateKeyBytes == nil {
		t.Fatal("Admin client private key not generated")
	}

	if adminClientCert.PublicKeyBytes == nil {
		t.Fatal("Admin client pubic key not generated")
	}
}
