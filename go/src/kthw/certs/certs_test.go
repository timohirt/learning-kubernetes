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

	if adminClientCert.PrivateKeyPath() != path.Join(adminClientCert.BaseDir, "admin-key.pem") {
		t.Fatalf("Private key path wrong. Should be '$baseDir/admin-key.pem'.")
	}

	if adminClientCert.PublicKeyBytes == nil {
		t.Fatal("Admin client pubic key not generated")
	}

	if adminClientCert.PublicKeyPath() != path.Join(adminClientCert.BaseDir, "admin.pem") {
		t.Fatalf("Public key path wrong. Should be '$baseDir/admin.pem'.")
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

func TestGenerateWriteReadEtcdCert(t *testing.T) {
	certGenerator, certsConfig := createCertGenerator(t)

	etcdCert, err := certGenerator.GenEtcdCertificate()
	helperFailIfErr(t, "Error while generating Etcd certificate", err)

	if etcdCert.PrivateKeyBytes == nil {
		t.Fatal("etcd private key not generated")
	}

	if etcdCert.PrivateKeyPath() != path.Join(etcdCert.BaseDir, "etcd-key.pem") {
		t.Fatalf("Private key path wrong. Should be '$baseDir/etcd-key.pem'.")
	}

	if etcdCert.PublicKeyBytes == nil {
		t.Fatal("etcd pubic key not generated")
	}

	if etcdCert.PublicKeyPath() != path.Join(etcdCert.BaseDir, "etcd.pem") {
		t.Fatalf("Public key path wrong. Should be '$baseDir/etcd.pem'.")
	}

	err = etcdCert.Write()
	if err != nil {
		t.Errorf("Error while writing etcd certificate: '%s'", err)
	}

	loadedCert, err := certs.LoadEtcdCert(certsConfig)
	if err != nil {
		t.Errorf("Error while loading etcd certificate: '%s'", err)
	}

	if !reflect.DeepEqual(loadedCert.PrivateKeyBytes, etcdCert.PrivateKeyBytes) {
		t.Error("Private key of loaded etcd cert differs from key generated and written to file.")
	}

	if !reflect.DeepEqual(loadedCert.PublicKeyBytes, etcdCert.PublicKeyBytes) {
		t.Error("Public key of loaded etcd cert differs from key generated and written to file.")
	}
}
