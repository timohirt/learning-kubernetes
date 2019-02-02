package cmd_test

import (
	"crypto/rsa"
	"io/ioutil"
	"kthw/cmd"
	"testing"

	"github.com/cloudflare/cfssl/helpers"
)

func TestGenerateCA(t *testing.T) {
	defaultCerts := cmd.DefaultCerts()
	tempDirName, err := ioutil.TempDir("", "InitCA")
	if err != nil {
		t.Fatalf("Error creating temp dir: %s", err)
	}

	defaultCerts.CABaseDir = tempDirName
	err = defaultCerts.InitCa()
	if err != nil {
		t.Fatalf("Error while generating CA: %s", err)
	}

	ca := defaultCerts.CA

	if ca.CertBytes == nil {
		t.Fatal("CA cert not generated")
	}
	if ca.KeyBytes == nil {
		t.Fatal("CA key not generated")
	}

	key, err := helpers.ParsePrivateKeyPEM(ca.KeyBytes)
	if err != nil {
		t.Fatalf("Error parsing generated private key: %s", err)
	}

	if key.(*rsa.PrivateKey).N.BitLen() != defaultCerts.CAKeySize {
		t.Fatalf("CA Private key lenght mismatch")
	}

	cert, err := helpers.ParseCertificatePEM(ca.CertBytes)
	if err != nil {
		t.Fatalf("Error parsing generated cert: %s", err)
	}

	if cert.PublicKey.(*rsa.PublicKey).N.BitLen() != defaultCerts.CAKeySize {
		t.Fatalf("CA Cert key lenght mismatch")
	}
}
