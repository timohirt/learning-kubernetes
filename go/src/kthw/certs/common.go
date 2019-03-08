package certs

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudflare/cfssl/csr"
)

var signingUsages = []string{"signing", "key encipherment", "server auth", "client auth"}

const (
	keySize = 2048
	keyAlgo = "rsa"

	caCN           = "Kubernetes"
	caKeyFileName  = "ca-key.pem"
	caCertFileName = "ca.pem"

	certsBaseDir = "pki"

	signingProfile = "kubernetes"
	signingExpiry  = "8760h"

	adminClientCN           = "admin"
	adminClientO            = "system:masters"
	adminClientKeyFileName  = "admin-key.pem"
	adminClientCertFileName = "admin.pem"

	etcdCN           = "etcd"
	etcdO            = "Kubernetes"
	etcdKeyFileName  = "etcd-key.pem"
	etcdCertFileName = "etcd.pem"
)

func certName(o string) csr.Name {
	return csr.Name{
		C:  "DE",
		L:  "Mainz",
		O:  o,
		OU: "Learning Kubernetes",
		ST: "RLP"}
}

func writeToFile(cert []byte, file string) error {
	var err error
	if _, statErr := os.Stat(file); statErr == nil {
		err = fmt.Errorf("Could not write certificate to already existing file %s", file)
	} else {
		err = ioutil.WriteFile(file, cert, 0644)
	}
	return err
}

func readFromFile(file string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("Error reading private key file '%s': %s", file, err)
	}
	return bytes, nil
}

func ensureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
