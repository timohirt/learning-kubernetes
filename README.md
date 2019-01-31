# Learning Kubernetes

This repository contains examples I used to learn more about Kubernetes. I use
this repository to take notes and write about what I learnt and what might be
useful to know about Kubernetes, its components and how to use it.

**Disclaimer**: As stated above, this repo is mainly for taking notes. It is
not meant to be a complete reference.

## Minikube

I primarily use [Minikube](https://kubernetes.io/docs/setup/minikube/) on my
MacBook. It is very easy to get a Kubernetes cluster up and running with VirtualBox.

```
$ minikube start --vm-driver=virtualbox
```

This command creates a VM and starts a sinlge node Kubernetes cluster. 

### Working with Swagger UI

Swagger UI is a great tool browse the API. Minikube can be configured to
enable it in the Kubernetes cluster.

```bash
$ minikube start --extra-config=apiserver.enable-swagger-ui=true
kubectl proxy --port=8080
```

Now open your browser and navigate to
[http://localhost:8080/swagger-ui/](http://localhost:8080/swagger-ui/).

## Making Docker Images Available in Minikube

Currently, I don't push Docker images to a public repository. Instead I
configure the `docker` command to use Docker on the Minicube VM to build an
image.

```bash
$ eval $(minikube docker-env)
```

After running this command all docker images build are available in Minikube
and can be used in pod specifications for example.

## Open a Shell in a Docker Container

Sometimes it is useful to open a shell in a Docker container to test something.
I use the `shell-demo` pod provided by the Kubernetes project. 

```bash
$ kubectl create -f https://k8s.io/examples/application/shell-demo.yaml
$ kubectl exec -it shell-demo -- bash
root@shell-demo:/#
```



