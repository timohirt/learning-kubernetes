package network

import (
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"strings"
)

// SetupWireguard generated wireguard config for each server and copies it using SCP.
func SetupWireguard(sshOperations sshconnect.SSHOperations, servers []*server.Config) error {
	wgConfs, _ := GenerateWireguardConf(servers)
	for _, hostConf := range wgConfs.WgHosts {
		hostIP := hostConf.PublicIP
		conf, _ := hostConf.generateServerConf()

		commands := &sshconnect.Commands{
			Commands: []sshconnect.Command{
				uploadConfigFile(hostIP, conf),
				openFirewall(hostIP),
				startDevice(hostIP)},
			LogOutput: true}

		sshOperations.RunCmds(commands)
	}
	return nil
}

func openFirewall(host string) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		Host:        host,
		CommandLine: "ufw allow in on wg0",
		Description: "Open firewall for private overlay network"}
}

func startDevice(host string) *sshconnect.ShellCommand {
	return &sshconnect.ShellCommand{
		Host:        host,
		CommandLine: "systemctl enable wg-quick@wg0 && systemctl restart wg-quick@wg0",
		Description: "Start wireguard device 'wg0'"}
}

func uploadConfigFile(host string, configFile string) *sshconnect.CopyFileCommand {
	return &sshconnect.CopyFileCommand{
		Host:        host,
		FileContent: strings.NewReader(configFile),
		FilePath:    "/etc/wireguard/wg0.conf",
		Description: "Upload wireguard config file of device 'wg0'"}
}
