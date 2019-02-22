package server

import (
	"kthw/cmd/hcloudclient"
	"kthw/cmd/infra/sshkey"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

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
  - linux-headers-generic
  - apt-transport-https
  - ca-certificates
  - curl
  - software-properties-common
  - [docker-ce, 18.06.1~ce~3-0~ubuntu]
runcmd:
  - [ sudo, ufw, allow, 22/tcp ]
  - [ sudo, ufw, allow, 51820/udp ]
  - [ sudo, ufw, enable ]
`

// Create creates a server in hcloud using the provided config. Public ip and
// root password are added to the conf and calling code is assumed to write the configuration.
func Create(config Config, client hcloudclient.HCloudOperations) (*Config, error) {
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

	serverCreated := client.Create(serverOpts)

	config.PublicIP = serverCreated.PublicIP
	config.RootPassword = serverCreated.RootPassword
	config.SSHPublicKeyID = sshKey.ID
	config.ID = serverCreated.ID

	return &config, nil
}
