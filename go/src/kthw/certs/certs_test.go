package certs_test

import (
	"kthw/certs"
	"path"
	"reflect"
	"testing"
)

func createCertGenerator(t *testing.T) (*certs.CertGenerator, certs.Config) {
	caCerts, certsDir := helperCreateDefaultCACerts(t)
	caCerts.InitCa()
	helperEnsureCaCertsInitialized(t, caCerts)
	certsConf := certs.Config{BaseDir: certsDir}
	certGenerator, err := certs.NewCertGenerator(caCerts, certsConf)
	helperFailIfErr(t, "Error while creating CertGenerator", err)

	return certGenerator, certsConf
}

func TestGenerateWriteReadAdminClientCert(t *testing.T) {
	certGenerator, certsConfig := createCertGenerator(t)

	adminClientCert, err := certGenerator.GenAdminClientCertificate()
	helperFailIfErr(t, "Error creating admin client certificate", err)

	if adminClientCert.PrivateKeyBytes == nil {
		t.Fatal("Admin client private key not generated")
	}

	if adminClientCert.PrivateKeyPath() != path.Join(adminClientCert.BaseDir, "admin.key") {
		t.Fatalf("Private key path wrong. Should be '$baseDir/admin.key'.")
	}

	if adminClientCert.PublicKeyBytes == nil {
		t.Fatal("Admin client pubic key not generated")
	}

	if adminClientCert.PublicKeyPath() != path.Join(adminClientCert.BaseDir, "admin.crt") {
		t.Fatalf("Public key path wrong. Should be '$baseDir/admin.crt'.")
	}

	err = adminClientCert.Write()
	if err != nil {
		t.Errorf("Error while writing admin client certificate: '%s'", err)
	}

	loadedCert, err := certs.LoadAdminClientCert(certsConfig)
	if err != nil {
		t.Errorf("Error while loading admin client certificate: '%s'", err)
	}

	if !reflect.DeepEqual(loadedCert.PrivateKeyBytes, adminClientCert.PrivateKeyBytes) {
		t.Error("Private key of loaded admin client cert differs from key generated and written to file.")
	}

	if !reflect.DeepEqual(loadedCert.PublicKeyBytes, adminClientCert.PublicKeyBytes) {
		t.Error("Public key of loaded admin client differs from key generated and written to file.")
	}
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
