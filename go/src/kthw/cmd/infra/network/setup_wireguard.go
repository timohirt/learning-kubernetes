package network

import (
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"strings"
)

// SetupWireguard generated wireguard config for each server and copies it using SCP.
func SetupWireguard(sshOperations sshconnect.SSHOperations, servers []server.Config) ([]server.Config, error) {
	wgConfs, _ := GenerateWireguardConf(servers)
	updatedServerConfig := []server.Config{}
	for _, hostConf := range wgConfs.WgHosts {
		sshIP := hostConf.PublicIP
		conf, _ := hostConf.generateServerConf()
		reader := strings.NewReader(conf)
		err := sshOperations.WriteReadOnlyFileTo(sshIP, reader, "/etc/wireguard/wg0.conf")
		if err != nil {
			return nil, err
		}
		_, err = sshOperations.RunCmd(sshIP, "systemctl enable wg-quick@wg0 && systemctl restart wg-quick@wg0")
		if err != nil {
			return nil, err
		}
		_, err = sshOperations.RunCmd(sshIP, "ufw allow in on wg0")
		if err != nil {
			return nil, err
		}
		updatedServerConfig = append(updatedServerConfig, hostConf.ServerConfig)
	}
	return updatedServerConfig, nil
}
