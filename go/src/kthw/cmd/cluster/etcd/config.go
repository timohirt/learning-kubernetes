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
  --listen-client-urls "http://{{.PrivateIP}}:2379,http://localhost:2379" \
  --advertise-client-urls "http://{{.PrivateIP}}:2379" \
  --initial-cluster "{{.NodeName}}=http://localhost:2380" \
  --heartbeat-interval 200 \
  --election-timeout 5000 \
  --cert-file /etc/kubernetes/pki/etcd.pem \
  --key-file /etc/kubernetes/pki/etcd-key.pem

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
