package etcd

import (
	"bytes"
	"fmt"
	"text/template"
)

type EtcdSystemdServiceParameters struct {
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
  --advertise-client-urls "http://{{.PrivateIP}}:2379"
  --initial-cluster "{{.NodeName}}=http://localhost:2380" \
  --heartbeat-interval 200 \
  --election-timeout 5000

[Install]
WantedBy=multi-user.target
`

// GenerateEtcdSystemdService generates a etcd systemd service.
func GenerateEtcdSystemdService(params EtcdSystemdServiceParameters) (string, error) {
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
