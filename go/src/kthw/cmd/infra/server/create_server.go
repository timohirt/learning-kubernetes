package server

import (
	"kthw/cmd/hcloudclient"
	"kthw/cmd/infra/sshkey"
	"kthw/cmd/sshconnect"
	"strings"

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
    kubernetes.list:
      source: "deb [arch=amd64] https://apt.kubernetes.io/ kubernetes-xenial main"
      keyid: BA07F4FB
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
  - kubelet
  - kubeadm
  - kubectl
runcmd:
  - [ sudo, ufw, allow, 22/tcp ]
  - [ sudo, ufw, allow, 51820/udp ]
  - [ sudo, ufw, enable ]
  - [ swapoff, -a ]
  - [ mkdir, -p, /etc/kubernetes/pki ]
  - [ apt-mark, hold, kubelet, kubeadm, kubectl, docker-ce ]
`

// Create creates a server in hcloud using the provided config. Public ip and
// root password are added to the conf and calling code is assumed to write the configuration.
func Create(config *Config, client hcloudclient.HCloudOperations) error {
	sshKeyFromConf, err := sshkey.ReadSSHPublicKeyFromConf()
	if err != nil {
		return err
	}

	serverType := &hcloud.ServerType{Name: config.ServerType}
	image := &hcloud.Image{Name: config.ImageName}
	location := &hcloud.Location{Name: config.LocationName}
	sshKey := &hcloud.SSHKey{ID: sshKeyFromConf.ID}
	startAfterCreate := true
	labels := make(map[string]string)
	labels["roles"] = strings.Join(config.Roles, ",")
	serverOpts := hcloud.ServerCreateOpts{
		Name:             config.Name,
		ServerType:       serverType,
		Image:            image,
		Location:         location,
		UserData:         basicCloudInit,
		SSHKeys:          []*hcloud.SSHKey{sshKey},
		StartAfterCreate: &startAfterCreate,
		Labels:           labels}

	serverCreated := client.Create(serverOpts)

	config.PublicIP = serverCreated.PublicIP
	config.RootPassword = serverCreated.RootPassword
	config.SSHPublicKeyID = sshKey.ID
	config.ID = serverCreated.ID

	return nil
}

// IsCloudInitCompleted tests if cloud-init already completed and returns 'true' if is did and otherwise 'false'
func IsCloudInitCompleted(host string, ssh sshconnect.SSHOperations) bool {
	command := checkBootFinishedFileExists(host)
	_, err := ssh.RunCmd(command, false)

	if err != nil {
		return false
	}
	return true
}

func checkBootFinishedFileExists(host string) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		Host:        host,
		CommandLine: "test -e /var/lib/cloud/instance/boot-finished",
		Description: "Check if cloud-init completed"}
}
