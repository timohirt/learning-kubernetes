package certs

type GeneratesCertsMock struct {
	ca                  *CA
	etcdCert            *EtcdCert
	etcdClientCert      *EtcdClientCert
	IsEtcdCertGenerated bool
	GeneratesCerts
}

// GetCA returns CA of mock.
func (g *GeneratesCertsMock) GetCA() *CA { return g.ca }

// GenEtcdCertificate returns etcd certificate of mock
func (g *GeneratesCertsMock) GenEtcdCertificate(hosts []string) (*EtcdCert, error) {
	g.IsEtcdCertGenerated = true
	return g.etcdCert, nil
}

// GenEtcdClientCertificate returns etcd client certificate of mock
func (g *GeneratesCertsMock) GenEtcdClientCertificate() (*EtcdClientCert, error) {
	return g.etcdClientCert, nil
}

// NewGeneratesCertsMock creates a new mock with dummy certs
func NewGeneratesCertsMock() *GeneratesCertsMock {
	ca := CA{
		CertBytes: []byte("CA_CERT"),
		KeyBytes:  []byte("CA_KEY")}

	etcdCert := EtcdCert{
		PrivateKeyBytes: []byte("ETCD_KEY"),
		PublicKeyBytes:  []byte("ETCD_CERT")}

	etcdClientCert := EtcdClientCert{
		PrivateKeyBytes: []byte("ETCD_KEY"),
		PublicKeyBytes:  []byte("ETCD_CERT")}

	return &GeneratesCertsMock{
		ca:                  &ca,
		etcdCert:            &etcdCert,
		etcdClientCert:      &etcdClientCert,
		IsEtcdCertGenerated: false}
}

type CertificateLoaderMock struct {
	etcdCert *EtcdClientCert
	ca       *CA
	CertificateLoader
}

// NewCertificateLoaderMock creates mock initialised with dummy certs.
func NewCertificateLoaderMock() *CertificateLoaderMock {
	return &CertificateLoaderMock{
		etcdCert: &EtcdClientCert{
			PrivateKeyBytes: []byte("ETCD_CLIENT_PRIVATE"),
			PublicKeyBytes:  []byte("ETCD_CLIENT_PUBLIC")},
		ca: &CA{
			CertBytes: []byte("CA_CERT"),
			KeyBytes:  []byte("CA_PRIVATE")},
	}
}

// LoadEtcdClientCert returns EtcdClientCert of mock.
func (c *CertificateLoaderMock) LoadEtcdClientCert() (*EtcdClientCert, error) {
	return c.etcdCert, nil
}

// LoadCA returns CA of mock.
func (c *CertificateLoaderMock) LoadCA() (*CA, error) {
	return c.ca, nil
}
