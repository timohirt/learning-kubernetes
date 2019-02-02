package cmd_test

import (
	"bytes"
	"crypto/rsa"
	"io/ioutil"
	"kthw/cmd"
	"testing"

	"github.com/cloudflare/cfssl/helpers"
)

func TestInitCA(t *testing.T) {
	defaultCaCerts, _ := helperCreateDefaultCACerts(t)
	err := defaultCaCerts.InitCa()
	helperFailIfErr(t, "Error while generating CA: %s", err)

	ca := defaultCaCerts.CA

	if ca.CertBytes == nil {
		t.Fatal("CA cert not generated")
	}
	if ca.KeyBytes == nil {
		t.Fatal("CA key not generated")
	}

	key, err := helpers.ParsePrivateKeyPEM(ca.KeyBytes)
	helperFailIfErr(t, "Error parsing generated private key: %s", err)

	if key.(*rsa.PrivateKey).N.BitLen() != defaultCaCerts.CAKeySize {
		t.Fatalf("CA Private key lenght mismatch")
	}

	cert, err := helpers.ParseCertificatePEM(ca.CertBytes)
	helperFailIfErr(t, "Error parsing generated cert: %s", err)

	if cert.PublicKey.(*rsa.PublicKey).N.BitLen() != defaultCaCerts.CAKeySize {
		t.Fatalf("CA Cert key lenght mismatch")
	}
}

func TestInitCACreatedFiles(t *testing.T) {
	defaultCaCerts, _ := helperCreateDefaultCACerts(t)
	err := defaultCaCerts.InitCa()
	helperFailIfErr(t, "Error while generating CA: %s", err)

	actualPrivateKeyBytes, err := ioutil.ReadFile(defaultCaCerts.CNPrivateKeyFile())
	helperFailIfErr(t, "Failed reading private key.", err)

	if !bytes.Equal(actualPrivateKeyBytes, defaultCaCerts.CA.KeyBytes) {
		t.Fatalf("Private key in *CACerts differs from key read from file.")
	}

	actualPublicKeyBytes, err := ioutil.ReadFile(defaultCaCerts.CNPublicKeyFile())
	helperFailIfErr(t, "Failed reading public key.", err)

	if !bytes.Equal(actualPublicKeyBytes, defaultCaCerts.CA.CertBytes) {
		t.Fatalf("Public key in *CACerts differs from key read from file.")
	}
}

func helperReadPEM(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

func helperCreateDefaultCACerts(t *testing.T) (*cmd.CACerts, string) {
	defaultCaCerts := cmd.DefaultCACerts()
	tempDirName, err := ioutil.TempDir("", "InitCA")
	helperFailIfErr(t, "Error creating temp dir: %s", err)

	defaultCaCerts.CABaseDir = tempDirName
	return defaultCaCerts, tempDirName
}

func helperFailIfErr(t *testing.T, message string, err error) {
	if err != nil {
		t.Fatalf(message, err)
	}
}
