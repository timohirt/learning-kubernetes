package certs_test

import (
	"kthw/certs"
	"path"
	"testing"
)

func createCertGenerator(t *testing.T) (*certs.CertGenerator, certs.Config) {
	caCerts, certsDir := helperCreateDefaultCACerts(t)
	caCerts.InitCa()
	helperEnsureCaCertsInitialized(t, caCerts)
	certsConf := certs.Config{BaseDir: certsDir}
	certGenerator, err := certs.NewCertGenerator(caCerts.CA, certsConf)
	helperFailIfErr(t, "Error while creating CertGenerator", err)

	return certGenerator, certsConf
}

func TestGenerateEtcdCert(t *testing.T) {
	certGenerator, _ := createCertGenerator(t)

	etcdCert, err := certGenerator.GenEtcdCertificate([]string{"localhost"})
	helperFailIfErr(t, "Error while generating Etcd certificate", err)

	if etcdCert.PrivateKeyBytes == nil {
		t.Fatal("etcd private key not generated")
	}

	if etcdCert.PublicKeyBytes == nil {
		t.Fatal("etcd pubic key not generated")
	}

	if etcdCert.PublicKeyPath() != path.Join(etcdCert.BaseDir, "etcd.crt") {
		t.Fatalf("Public key path wrong. Should be '$baseDir/etcd.crt'.")
	}

	if etcdCert.PrivateKeyPath() != path.Join(etcdCert.BaseDir, "etcd.key") {
		t.Fatalf("Private key path wrong. Should be '$baseDir/etcd.key'.")
	}
}

func TestGenerateEtcdClientCert(t *testing.T) {
	certGenerator, _ := createCertGenerator(t)

	etcdClientCert, err := certGenerator.GenEtcdClientCertificate()
	helperFailIfErr(t, "Error while generating Etcd client certificate", err)

	if etcdClientCert.PrivateKeyBytes == nil {
		t.Fatal("etcd client private key not generated")
	}

	if etcdClientCert.PublicKeyBytes == nil {
		t.Fatal("etcd client pubic key not generated")
	}

	if etcdClientCert.PublicKeyPath() != path.Join(etcdClientCert.BaseDir, "etcd-client.crt") {
		t.Fatalf("Public key path wrong. Should be '$baseDir/etcd-client.crt'.")
	}

	if etcdClientCert.PrivateKeyPath() != path.Join(etcdClientCert.BaseDir, "etcd-client.key") {
		t.Fatalf("Private key path wrong. Should be '$baseDir/etcd-client.key'.")
	}
}
