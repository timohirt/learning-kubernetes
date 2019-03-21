package kube

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"kthw/cmd/infra/server"
	"kthw/cmd/sshconnect"
	"strings"
)

type AddOn interface {
	Description() string
	getCommands(hostConfig *server.Config) []sshconnect.Command
}

// CalicoNetworkingAddOn installs Calico networking to a K8s cluster
type CalicoNetworkingAddOn struct {
	rbacManifest   string
	calicoManifest string
	AddOn
}

func NewCalicoNetworkingAddOn() *CalicoNetworkingAddOn {
	return &CalicoNetworkingAddOn{
		rbacManifest:   "https://docs.projectcalico.org/v3.3/getting-started/kubernetes/installation/hosted/rbac-kdd.yaml",
		calicoManifest: "https://docs.projectcalico.org/v3.3/getting-started/kubernetes/installation/hosted/kubernetes-datastore/calico-networking/1.7/calico.yaml",
	}
}

func (c *CalicoNetworkingAddOn) Description() string { return "Install Calico networking" }

func (c *CalicoNetworkingAddOn) getCommands(hostConfig *server.Config) []sshconnect.Command {
	podNetwork := strings.Replace(podNetworkCIDR, "/", "\\/", -1)
	downloadManifests := fmt.Sprintf("curl %s -o /tmp/rbac-kdd.yaml && curl %s -o /tmp/calico.yaml", c.rbacManifest, c.calicoManifest)
	changePodNetwork := fmt.Sprintf("sed -i 's/192.168.0.0\\/16/%s/g' /tmp/calico.yaml", podNetwork)
	installCalico := "kubectl apply -f /tmp/rbac-kdd.yaml && kubectl apply -f /tmp/calico.yaml"
	return []sshconnect.Command{
		&sshconnect.ShellCommand{
			CommandLine: fmt.Sprintf("%s && %s && %s", downloadManifests, changePodNetwork, installCalico),
			Host:        hostConfig.PublicIP,
			Description: c.Description()}}
}

// KubernetesDashboardAddOn installs Kubernetes dashboard to a cluster and creates an admin user.
type KubernetesDashboardAddOn struct {
	dashboardManifest string
	AddOn
}

func NewKubernetesDashboardAddOn() *KubernetesDashboardAddOn {
	return &KubernetesDashboardAddOn{
		dashboardManifest: "https://raw.githubusercontent.com/kubernetes/dashboard/master/aio/deploy/recommended/kubernetes-dashboard.yaml",
	}
}

func (c *KubernetesDashboardAddOn) Description() string { return "Install Kubernetes dashboard" }

func (c *KubernetesDashboardAddOn) getCommands(hostConfig *server.Config) []sshconnect.Command {
	fileContent, err := ioutil.ReadFile("manifests/dashboard-admin.yaml")
	if err != nil {
		fmt.Printf("Error while reading dashboard admin manifest: %s", err)
		return []sshconnect.Command{}
	}
	return []sshconnect.Command{
		&sshconnect.CopyFileCommand{
			Host:        hostConfig.PublicIP,
			FileContent: bytes.NewReader(fileContent),
			FilePath:    "/tmp/dashboard-admin.yaml",
			Description: "Copy dashboard-admin-yaml to cluster"},
		&sshconnect.ShellCommand{
			CommandLine: fmt.Sprintf("kubectl apply -f /tmp/dashboard-admin.yaml"),
			Host:        hostConfig.PublicIP,
			Description: c.Description()},
		&sshconnect.ShellCommand{
			CommandLine: fmt.Sprintf("kubectl apply -f %s", c.dashboardManifest),
			Host:        hostConfig.PublicIP,
			Description: c.Description()}}
}
