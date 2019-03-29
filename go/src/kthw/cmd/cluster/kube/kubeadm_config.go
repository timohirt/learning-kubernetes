package kube

import (
	"bytes"
	"fmt"
	"html/template"
	"kthw/cmd/infra/server"
)

var controllerConfigTemplate = `
apiVersion: kubeadm.k8s.io/v1beta1
kind: InitConfiguration
bootstrapTokens:
- groups:
  - system:bootstrappers:kubeadm:default-node-token
  token: abcdef.0123456789abcdef
  ttl: 24h0m0s
  usages:
  - signing
  - authentication
localAPIEndpoint:
  advertiseAddress: {{.PublicIP}}
  bindPort: 6443
nodeRegistration:
  name: {{.NodeName}}
  taints:
  - effect: NoSchedule
    key: node-role.kubernetes.io/master
---
apiVersion: kubeadm.k8s.io/v1beta1
kind: ClusterConfiguration
clusterName: kubernetes
etcd:
  external:
    endpoints:{{range .EtcdNodes}}
    - {{.EndpointURL}}
    {{end}}caFile: /etc/kubernetes/pki/ca.crt
    certFile: /etc/kubernetes/pki/etcd-client.crt
    keyFile: /etc/kubernetes/pki/etcd-client.key
kubernetesVersion: v1.14.0
networking:
  dnsDomain: cluster.local
  podSubnet: "{{.PodNetworkCIDR}}"
  serviceSubnet: "10.96.0.0/12"
`

const podNetworkCIDR = "10.100.0.0/16"

type KubeAdmParams struct {
	PrivateIP      string
	PublicIP       string
	NodeName       string
	EtcdNodes      []EtcdNode
	PodNetworkCIDR string
}

type EtcdNode struct {
	EndpointURL string
}

func NewKubeAdmParams(hostConfig *server.Config, etcdNodes []EtcdNode) KubeAdmParams {
	return KubeAdmParams{
		PrivateIP:      hostConfig.PrivateIP,
		PublicIP:       hostConfig.PublicIP,
		NodeName:       hostConfig.Name,
		EtcdNodes:      etcdNodes,
		PodNetworkCIDR: podNetworkCIDR}
}

// GenerateKubeadmControllerConfig generates kubeadm controller config file
func GenerateKubeadmControllerConfig(params KubeAdmParams) (string, error) {
	tmpl, err := template.New("controller-config").Parse(controllerConfigTemplate)
	if err != nil {
		return "", err
	}

	var confBuffer bytes.Buffer
	err = tmpl.ExecuteTemplate(&confBuffer, "controller-config", params)
	if err != nil {
		fmt.Println("Failed generating kubernetes controller config.")
		return "", err
	}

	return confBuffer.String(), nil
}
