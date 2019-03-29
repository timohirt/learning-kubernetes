# Build

```bash
make deps
make build
```

# Creating a K8s Cluster

```bash
$ ./kthw project new p1 ~/.ssh/id_ed25519.pub --apiToken <token>
Initialised project.yaml with defaults.
Added SSH key to config.
SSH key created in hcloud.
CA private and public keys generated and stored in pkiInitialised PKI infrastructure.

$ ./kthw project add-server controller-1 etcd,controller
Server controller-1 successfully added to config.

$ ./kthw project add-server controller-2 worker
Server controller-2 successfully added to config.

$ ./kthw install k8s-non-ha --apiToken <token>
Creating server controller-1 at Hetzner cloud
Creating server controller-2 at Hetzner cloud
Waiting for controller-2 to complete cloud-init
Waiting for controller-1 to complete cloud-init
Check if cloud-init completed
Check if cloud-init completed
```
