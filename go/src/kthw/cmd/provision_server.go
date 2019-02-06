package cmd

import (
	"kthw/cmd/config"
	"kthw/cmd/hcloudclient"
	"kthw/cmd/sshkey"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// CreateServerOpts are used to create a server in hcloud
type CreateServerOpts struct {
	name         string
	serverType   string
	imageName    string
	locationName string
}

var basicCloudInit = `#cloud-config
apt:
  preserve_sources_list: true
  sources:
    docker-ppa.list:
      source: "deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable"
      keyid: 0EBFCD88 
    wireguard-ppa.list:
      source: "ppa:wireguard/wireguard"
      keyid: 504A1A25
apt_update: true
apt_upgrade: true
apt_reboot_if_required: true 
packages:
  - wireguard
  - apt-transport-https
  - ca-certificates
  - curl
  - software-properties-common
  - [docker-ce, 18.06.1~ce~3-0~ubuntu]
`

// CreateServer creates a server in hcloud using the provided config. Public ip and
// root password are added to the conf and calling code is assumed to write the configuration.
func createServer(config config.ServerConfig, client hcloudclient.HCloudOperations) (*config.ServerConfig, error) {
	sshKeyFromConf, err := sshkey.ReadSSHPublicKeyFromConf()
	if err != nil {
		return nil, err
	}

	serverType := &hcloud.ServerType{Name: config.ServerType}
	image := &hcloud.Image{Name: config.ImageName}
	location := &hcloud.Location{Name: config.LocationName}
	sshKey := &hcloud.SSHKey{ID: sshKeyFromConf.ID}
	startAfterCreate := true
	serverOpts := hcloud.ServerCreateOpts{
		Name:             config.Name,
		ServerType:       serverType,
		Image:            image,
		Location:         location,
		UserData:         basicCloudInit,
		SSHKeys:          []*hcloud.SSHKey{sshKey},
		StartAfterCreate: &startAfterCreate}

	serverCreated := client.CreateServer(serverOpts)

	config.PublicIP = serverCreated.PublicIP
	config.RootPassword = serverCreated.RootPassword
	config.SSHPublicKeyID = sshKey.ID

	return &config, nil
}
