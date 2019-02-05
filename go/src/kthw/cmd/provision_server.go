package cmd

import (
	"fmt"

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
func CreateServer(config ServerConfig, client HCloudOperations) (*ServerConfig, error) {
	serverType := &hcloud.ServerType{Name: config.ServerType}
	image := &hcloud.Image{Name: config.ImageName}
	location := &hcloud.Location{Name: config.LocationName}
	startAfterCreate := true
	serverOpts := hcloud.ServerCreateOpts{
		Name:             config.Name,
		ServerType:       serverType,
		Image:            image,
		Location:         location,
		UserData:         basicCloudInit,
		StartAfterCreate: &startAfterCreate}

	serverCreated, err := client.CreateServer(serverOpts)
	if err != nil {
		return nil, fmt.Errorf("Error creating server %s. %s", config.Name, err)
	}

	config.PublicIP = serverCreated.PublicIP
	config.RootPassword = serverCreated.RootPassword

	return &config, nil
}