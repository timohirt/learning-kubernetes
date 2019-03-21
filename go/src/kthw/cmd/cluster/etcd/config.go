package etcd

import (
	"bytes"
	"fmt"
	"text/template"
)

// SystemdServiceParameters parameters used to generate etcs systemd service.
type SystemdServiceParameters struct {
	PrivateIP string
	NodeName  string
}

var etcdSystemdService = `[Unit]
Description=etcd

[Service]
Type=notify
Restart=always
RestartSec=5
TimeoutStartSec=0
StartLimitInterval=0
ExecStart=/usr/local/bin/etcd \
  --name {{.NodeName}} \
  --data-dir /var/lib/etcd \
  --listen-client-urls "https://{{.PrivateIP}}:2379,https://localhost:2379" \
  --advertise-client-urls "https://{{.PrivateIP}}:2379" \
  --initial-cluster "{{.NodeName}}=https://{{.PrivateIP}}:2380" \
  --initial-advertise-peer-urls "https://{{.PrivateIP}}:2380" \
  --listen-peer-urls "https://{{.PrivateIP}}:2380" \
  --heartbeat-interval 200 \
  --election-timeout 5000 \
  --trusted-ca-file /etc/etcd/pki/ca.crt \
  --cert-file /etc/etcd/pki/etcd.crt \
  --key-file /etc/etcd/pki/etcd.key \
  --client-cert-auth \
  --peer-trusted-ca-file /etc/etcd/pki/ca.crt \
  --peer-cert-file /etc/etcd/pki/etcd.crt \
  --peer-key-file /etc/etcd/pki/etcd.key \
  --peer-client-cert-auth

[Install]
WantedBy=multi-user.target
`

// GenerateSystemdService generates a etcd systemd service.
func GenerateSystemdService(params SystemdServiceParameters) (string, error) {
	tmpl, err := template.New("etcd-service").Parse(etcdSystemdService)
	if err != nil {
		return "", err
	}

	var confBuffer bytes.Buffer
	err = tmpl.ExecuteTemplate(&confBuffer, "etcd-service", params)
	if err != nil {
		fmt.Println("Failed generating etcd systemd service.")
		return "", err
	}

	return confBuffer.String(), nil
}
