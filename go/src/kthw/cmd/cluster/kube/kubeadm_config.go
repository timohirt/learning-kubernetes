package kube

var configTemplate = `
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
  advertiseAddress: {{.PrivateIP}}
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
	endpoints:
	  {{range .EtcdNodes}}
	  - {{.EndpointURL}}
	  {{end}}
kubernetesVersion: v1.13.0
networking:
  dnsDomain: cluster.local
  podSubnet: "10.100.0.0/16"
  serviceSubnet: 10.96.0.0/12
`
