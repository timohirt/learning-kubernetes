package certs_test

import (
	"kthw/certs"
	"testing"
)

func TestInitDefaultConfigAndReadConfig(t *testing.T) {
	certs.InitDefaultConfig()

	conf := certs.ReadConfig()
	if conf.BaseDir != "pki" {
		t.Errorf("Expected base dir is 'pki', but was %s", conf.BaseDir)
	}
}
