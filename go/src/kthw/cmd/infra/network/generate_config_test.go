package network_test

import (
	"kthw/cmd/infra/network"
	"kthw/cmd/infra/server"
	"reflect"
	"testing"
)

func TestGenerateConfigWithOneHost(t *testing.T) {
	serverConfigs := []server.Config{server.Config{ID: 1, PublicIP: "192.168.1.1"}}

	wgConfs, err := network.GenerateWireguardConf(serverConfigs)
	if err != nil {
		t.Fatalf("Error while generating host confs: %s", err)
	}
	host := wgConfs.WgHosts[0]

	ensurePrivateIPIsSet(host.ServerConfig, t)
	ensureHostFieldsAreFilled(host, t)
}

func TestGenerateConfigWithThreeHosts(t *testing.T) {
	hostConfigs := []server.Config{
		server.Config{ID: 1, PublicIP: "192.168.1.1"},
		server.Config{ID: 2, PublicIP: "192.168.1.2"},
		server.Config{ID: 3, PublicIP: "192.168.1.3"}}

	wgConf, err := network.GenerateWireguardConf(hostConfigs)
	if err != nil {
		t.Fatalf("Error while generating host confs: %s", err)
	}

	if len(wgConf.WgHosts) != 3 {
		t.Fatalf("Expected two hosts, but found %d", len(wgConf.WgHosts))
	}
	host1 := wgConf.WgHosts[0]
	host2 := wgConf.WgHosts[1]
	host3 := wgConf.WgHosts[2]
	ensurePrivateIPIsSet(host1.ServerConfig, t)
	ensurePrivateIPIsSet(host2.ServerConfig, t)
	ensurePrivateIPIsSet(host3.ServerConfig, t)
	ensureHostFieldsAreFilled(host1, t)
	ensureHostFieldsAreFilled(host2, t)
	ensureHostFieldsAreFilled(host3, t)

	ensureHostHasPeer(host1, []*network.Host{host2, host3}, t)
	ensureHostHasPeer(host2, []*network.Host{host1, host3}, t)
	ensureHostHasPeer(host3, []*network.Host{host1, host2}, t)
}

func ensureHostFieldsAreFilled(host *network.Host, t *testing.T) {
	if host.PrivateIP == "" {
		t.Fatalf("Expected PrivateIP to be non-empty string, but it is empty")
	}
	if host.PublicIP == "" {
		t.Fatalf("Expected PublicIP to be non-empty string, but it is empty")
	}
	if host.PublicKey == "" {
		t.Fatalf("Expected PublicKey to be non-empty string, but it is empty")
	}
	if host.PrivateKey == "" {
		t.Fatalf("Expected PrivateKey to be non-empty string, but it is empty")
	}
}

func ensureHostHasPeer(host *network.Host, peerHosts []*network.Host, t *testing.T) {
	var expectedPeers []network.Peer
	for _, peerHost := range peerHosts {
		expectedPeer := network.Peer{
			PublicKey:  peerHost.PublicKey,
			AllowedIPs: peerHost.PrivateIP,
			Endpoint:   peerHost.PublicIP}
		expectedPeers = append(expectedPeers, expectedPeer)

	}

	if len(host.Peers) != len(expectedPeers) {
		t.Errorf("Expected '%d' peer(s), but host has '%d' peers.", len(expectedPeers), len(host.Peers))
	}

	for _, expectedPeer := range expectedPeers {
		hasPeer := peerExists(host.Peers, expectedPeer)
		if !hasPeer {
			t.Errorf("Host %d was expected to have peer %s, but it hasn't.\n", host.ServerConfig.ID, expectedPeer.PublicKey)
		}
	}
}

func ensurePrivateIPIsSet(conf server.Config, t *testing.T) {
	if conf.PrivateIP == "" {
		t.Fatalf("PrivateIP of server %d not set", conf.ID)
	}
}

func peerExists(peers []network.Peer, peer network.Peer) bool {
	for _, p := range peers {
		if reflect.DeepEqual(p, peer) {
			return true
		}
	}
	return false
}
