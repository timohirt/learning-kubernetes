package certs_test

import (
	"kthw/certs"
	"testing"
)

func TestReadWriteCertConfig(t *testing.T) {
	certs.WriteDefaultConfig()

	conf := certs.ReadConfig()
	if conf.BaseDir != "pki" {
		t.Errorf("Expected base dir is 'pki', but was %s", conf.BaseDir)
	}
}
