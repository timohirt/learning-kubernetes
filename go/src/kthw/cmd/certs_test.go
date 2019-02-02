package cmd_test

import (
	"bytes"
	"crypto/rsa"
	"io/ioutil"
	"kthw/cmd"
	"log"
	"testing"

	"github.com/cloudflare/cfssl/helpers"
)

func TestInitCA(t *testing.T) {
	defaultCaCerts := helperCreateDefaultCACerts(t)
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
	defaultCaCerts := helperCreateDefaultCACerts(t)
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

func TestInitCANotOverrideExistingPrivateKey(t *testing.T) {
	defaultCaCerts := helperCreateDefaultCACerts(t)

	err := helperCreateFile(defaultCaCerts.CNPrivateKeyFile())
	if err != nil {
		t.Fatalf("Error setting up file for test %s", defaultCaCerts.CNPrivateKeyFile())
	}

	err = defaultCaCerts.InitCa()
	if err == nil {
		t.Fatalf("Existing private key file in output directory overridden!")
	}
}

func TestInitCANotOverrideExistingPublicKey(t *testing.T) {
	defaultCaCerts := helperCreateDefaultCACerts(t)

	helperCreateFile(defaultCaCerts.CNPublicKeyFile())

	err := defaultCaCerts.InitCa()
	if err == nil {
		t.Fatalf("Existing public key file in output directory overridden!")
	}
}

func helperCreateDefaultCACerts(t *testing.T) *cmd.CACerts {
	defaultCaCerts := cmd.DefaultCACerts()
	tempDirName, err := ioutil.TempDir("", "InitCA")
	helperFailIfErr(t, "Error creating temp dir: %s", err)

	defaultCaCerts.CABaseDir = tempDirName
	return defaultCaCerts
}

func helperReadPEM(file string) ([]byte, error) {
	return ioutil.ReadFile(file)
}

func helperCreateFile(file string) error {
	log.Printf("Dummy file: %s", file)
	var irrelevantContent [20]byte
	copy(irrelevantContent[:], "irrelevant")

	return ioutil.WriteFile(file, irrelevantContent[:], 0644)
}

func helperFailIfErr(t *testing.T, message string, err error) {
	if err != nil {
		t.Fatalf(message, err)
	}
}
