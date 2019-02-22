package network

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"kthw/cmd/infra/server"
	"os"
	"text/template"

	"golang.org/x/crypto/curve25519"
)

var serverInterfaceTemplate = `[Interface]
PrivateKey = {{.PrivateKey}}
ListenPort = 51820
Address = {{.PrivateIP}}
{{range .Peers}}

[Peer]
PublicKey = {{.PublicKey}}
Endpoint = {{.Endpoint}}:51820
AllowedIPs = {{.AllowedIPs}}
{{end}}
`

// Host holds IPs, public and private key used to connect servers with wireguard.
type Host struct {
	PublicIP     string
	PrivateIP    string
	PrivateKey   string
	PublicKey    string
	Peers        []Peer
	ServerConfig server.Config
}

// ToPeer generates a Peer from the fields of a Host
func (h *Host) ToPeer() Peer {
	if h.PublicKey == "" || h.PrivateIP == "" || h.PublicIP == "" {
		fmt.Printf("Error converting host to peer. Private key (%s), private ip (%s) and public key (%s) must be non-zero", h.PrivateKey, h.PrivateIP, h.PublicKey)
		os.Exit(1)
	}
	return Peer{
		PublicKey:  h.PublicKey,
		AllowedIPs: h.PrivateIP,
		Endpoint:   h.PublicIP}
}

func (h *Host) generateServerConf() (string, error) {
	tmpl, err := template.New("interface").Parse(serverInterfaceTemplate)
	if err != nil {
		return "", err
	}

	var confBuffer bytes.Buffer
	tmpl.ExecuteTemplate(&confBuffer, "interface", h)

	renderedConfig := confBuffer.String()
	return renderedConfig, nil
}

// Peer configuration of wireguard configuration.
type Peer struct {
	PublicKey  string
	AllowedIPs string
	Endpoint   string
}

// WgConf contains the wireguard configuration for each host.
type WgConf struct {
	WgHosts []*Host
}

// GenerateWireguardConf generates wireguard configuration for all servers passed on.
func GenerateWireguardConf(servers []server.Config) (*WgConf, error) {
	hosts := genAndAddKeys(servers)

	peers := newPeers(hosts)
	for _, host := range hosts {
		allOtherPeers := peers.selectAllExcept(host)
		host.Peers = allOtherPeers
	}

	conf := &WgConf{WgHosts: hosts}

	return conf, nil
}

type peers struct {
	all []Peer
}

func (p *peers) selectAllExcept(h *Host) []Peer {
	var selected []Peer
	for _, peer := range p.all {
		if peer.PublicKey != h.PublicKey {
			selected = append(selected, peer)
		}
	}
	return selected
}

func newPeers(hosts []*Host) peers {
	allPeers := make([]Peer, len(hosts))
	for cnt, host := range hosts {
		allPeers[cnt] = host.ToPeer()
	}
	return peers{all: allPeers}
}

type internalIPGenerator struct {
	IPCount int
}

func (i *internalIPGenerator) nextIP() string {
	i.IPCount = i.IPCount + 1
	return fmt.Sprintf("10.0.0.%d", i.IPCount)
}

func genAndAddKeys(serverConfigs []server.Config) []*Host {
	ipGen := internalIPGenerator{}
	results := make([]*Host, len(serverConfigs))
	for count, conf := range serverConfigs {
		keyPair, err := generateKeyPair()
		if err != nil {
			fmt.Printf("Error generating private key: %s", err)
			os.Exit(1)
		}

		conf.PrivateIP = ipGen.nextIP()
		host := Host{
			PublicIP:     conf.PublicIP,
			PrivateIP:    conf.PrivateIP,
			PrivateKey:   keyPair.private,
			PublicKey:    keyPair.public,
			ServerConfig: conf}
		results[count] = &host
	}
	return results
}

type keyPair struct {
	private string
	public  string
}

// generateKeyPair create a key-pair used to instantiate a wireguard connection
// Code copied and adapted from https://github.com/WireGuard/wireguard-go/blob/1c025570139f614f2083b935e2c58d5dbf199c2f/noise-helpers.go
func generateKeyPair() (keyPair, error) {
	var publicKey [32]byte
	var privateKey [32]byte
	_, err := rand.Read(privateKey[:])
	if err != nil {
		return keyPair{}, fmt.Errorf("unable to generate a private key: %v", err)
	}

	privateKey[0] &= 248
	privateKey[31] = (privateKey[31] & 127) | 64

	curve25519.ScalarBaseMult(&publicKey, &privateKey)

	return keyPair{
		private: base64.StdEncoding.EncodeToString(privateKey[:]),
		public:  base64.StdEncoding.EncodeToString(publicKey[:]),
	}, nil
}
