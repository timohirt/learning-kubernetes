# Services

A Kubernetes Service is a single, constant entrypoint to a group of pods. A
client can send HTTP requests to a Service which routes them to pods. It
balances the load between all selected pods. They are very much like typical
load balancers, and depending on the selected type, they represent an actual
load balancer in the underlyinfg infrastructure.

Services allow for easy horizontal scaling, as new pods are automatically
identified by at traffic routed to them. Additionally, using a Service instead
of connecting directly to a pod allows to move pods in a cluster at any time,
because the IP of the service stays the same.

## When to use Services

Imagine a typical web application. You might have several frontend pod. An
external client wants to connect to one of these pods. A Service provides a
single IP this client can use, instead of connecting to individual pods and
handling connection errors because they disappeared.

There might also be a database running in a pod which is not accessible for
external clients. The frondend pods need to connect to the database, but the
database pod will eventually be rescheduled and moved to another Kubernetes
Node and change its IP address. You don't want to reconfigure the frontends then.
Instead you would connect to a Service providing a static IP and taking care of
routing the traffic to the database, wherever the pod might live in the
cluster.

## Creating Services

Take a look at `simple-http.service.yaml` for the full resource. It defines a
ReplicaSet which created pods with `app: simple-httpr-.svc` labels. This is the
yaml of a Service.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: simple-http-svc
spec:
  ports:
  - port: 30001
    targetPort: 30000
  selector:
    app: simple-http-svc
```

It creates a Service listening on Port 30001, routing all connections to port
30000. The pods are identified using a `labelSelector`. In this case it routes
to all pods having a `app=simple-http-svd` label.

```bash
$ kubectl get svc
NAME              TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)     AGE
kubernetes        ClusterIP   10.96.0.1      <none>        443/TCP     4d
simple-http-svc   ClusterIP   10.97.201.12   <none>        30001/TCP   44m
```

The listing above shows that the IP assigned to the `simple-http-svc` Service
is a ClusterIp. This means that the service is only reachable from within the
cluster. We will later see how to make a service available for external clients.

## Discoverying Services

By default you can choose between using Kubernetes DNS or environment
variables. The first option creates a DNS record for each Service, like
`simple-http-svc.default.svc.cluster.local`. The first part is the name of the
Service, the second the namespace in Kubernetes (in this case default), and the
rest is static for services.

Let's take a look at the second option. When a Docker Image in a pod is
started, several env variables are passed to it. To shwo them, we create a new
pod with a shell-demo Docker container and the start an interactive shell.

```bash
$ kubectl exec -it shell-demo -- env
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOSTNAME=shell-demo
SIMPLE_HTTP_SVC_SERVICE_PORT=30001
SIMPLE_HTTP_SVC_SERVICE_HOST=10.97.201.12
...
KUBERNETES_SERVICE_HOST=10.96.0.1
KUBERNETES_SERVICE_PORT=443
...
```

The listing above shows the env variables `SIMPLE_HTTP_SERVICE_HOST` and
`SIMPLE_HTTP_SERVICE_PORT`, which can be used in a Docker container to connect
to the just created Service. Let's try this.

```bash
$ kubectl exec -it shell-demo -- bash
root@shell-demo:/# curl "http://$SIMPLE_HTTP_SVC_SERVICE_HOST:$SIMPLE_HTTP_SVC_SERVICE_PORT/health"
{"status":"OK","hostname":"simple-http-svc-sdbhp"}
root@shell-demo:/# curl "http://$SIMPLE_HTTP_SVC_SERVICE_HOST:$SIMPLE_HTTP_SVC_SERVICE_PORT/health"
{"status":"OK","hostname":"simple-http-svc-sdbhp"}
root@shell-demo:/# curl "http://$SIMPLE_HTTP_SVC_SERVICE_HOST:$SIMPLE_HTTP_SVC_SERVICE_PORT/health"
{"status":"OK","hostname":"simple-http-svc-sdbhp"}
root@shell-demo:/# curl "http://$SIMPLE_HTTP_SVC_SERVICE_HOST:$SIMPLE_HTTP_SVC_SERVICE_PORT/health"
{"status":"OK","hostname":"simple-http-svc-sdbhp"}
root@shell-demo:/# curl "http://$SIMPLE_HTTP_SVC_SERVICE_HOST:$SIMPLE_HTTP_SVC_SERVICE_PORT/health"
{"status":"OK","hostname":"simple-http-svc-8nn7n"}
root@shell-demo:/# curl "http://$SIMPLE_HTTP_SVC_SERVICE_HOST:$SIMPLE_HTTP_SVC_SERVICE_PORT/health"
{"status":"OK","hostname":"simple-http-svc-sdbhp"}
root@shell-demo:/# curl "http://$SIMPLE_HTTP_SVC_SERVICE_HOST:$SIMPLE_HTTP_SVC_SERVICE_PORT/health"
{"status":"OK","hostname":"simple-http-svc-8nn7n"}
root@shell-demo:/#
```

It works. Take a look at the `hostname` in the response. It even balanced the
load accross both pods currently running.

## Using Service to Connect to Services Outside the Cluster

You might want to connect pods to Services/Applications running outside the
Cluster. Using a Kubernetes Service instead of directly connecting to it
provides load balancing and service discovery. I won't go into detail, but will
outline how it works.

```bash
$ kubectl describe svc simple-http-svc
Name:              simple-http-svc
Namespace:         default
Labels:            <none>
Annotations:       <none>
Selector:          app=simple-http-svc
Type:              ClusterIP
IP:                10.97.201.12
Port:              <unset>  30001/TCP
TargetPort:        30000/TCP
Endpoints:         172.17.0.6:30000,172.17.0.8:30000
Session Affinity:  None
Events:            <none>
```

The listing above shows the details of the Services we created before. The
Endpoints resource lists the IPs and port of all pods matched by the pod
selector. Kubernetes Service created this automatically for you.

If no selector is specified, no Endoints resource is generated. Instead you can
create one manually, which must have the same name as the Service, and which
lists IPs or FQDN of external services, the Kubernetes Service routes
connections to.

## Exposing Services to External Clients

Until now our Service used a cluster ip and was thus only reachabel from within
a clustern. In oder to allow external clients to connect to a Service, there
are three options:

- Use service type `NodePort`. Opens a port on all nodes of a cluster.
  Connections to node `IP address:NodePort` are routed to the Service.
- User service type `LoadBalancer`. Depending on the environment your cluster
  runs in, a dedicated load balancer in your infrastructure is created and
  routes traffic to the pods. (Not supported in Minikube)
- Use an Ingress resource. It operates at the HTTP level, like AWS Application
  Load Balancer, and often routes traffic based on hostname header or path.. 

### Use the NodePort

Let's create a new service.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: simple-http-svc-nodeport
spec:
  type: NodePort
  ports:
  - port: 30001
    targetPort: 30000
    nodePort: 30002
  selector:
    app: simple-http-svc
```

The name of the service changed and there is now a `type: NodePort` set in the
`spec`. When you create this resource, port 30001 is opened on each node in the
Kubernetes cluster. (You may have to open port in your firewall to be able to
establish a connection).

```bash
$ kubectl create -f simple-http-service-nodeport.yaml
service "simple-http-svc-nodeport" created
replicaset.apps "simple-http-svc" created
```

The yaml in the repo contains a Service with `type` NodePort and the pod we
already know from other examples. Let's take a look at the details of the
Service.

```bash
$ kubectl describe svc simple-http-svc-nodeport
Name:                     simple-http-svc-nodeport
Namespace:                default
Labels:                   <none>
Annotations:              <none>
Selector:                 app=simple-http-svc
Type:                     NodePort
IP:                       10.110.56.154
Port:                     <unset>  30001/TCP
TargetPort:               30000/TCP
NodePort:                 <unset>  30002/TCP
Endpoints:                172.17.0.4:30000,172.17.0.6:30000
Session Affinity:         None
External Traffic Policy:  Cluster
Events:                   <none>
```

Now, you can send requests to port 30002 of the cluster node.

```bash
$ http http://192.168.99.104:30002/health
HTTP/1.1 200 OK
Content-Length: 51
Content-Type: text/plain; charset=utf-8
Date: Mon, 28 Jan 2019 08:21:40 GMT

{
    "hostname": "simple-http-svc-v85zm",
    "status": "OK"
}
```

Hint: Use `minikube service simple-http-svc-nodeport --url` to get this url.

Under the hood, every incoming connection to port 30002 is routed to a randomly
selected pod. It might be on this node or on another.

### Using LoadBalancer

If you change the `type` from `NodePort` to `LoadBalancer`, a load balancer in
the infrastructure is provisioned. Services are then available by the IP of
this load balancer. It still opens a port on each node, and the load balancer
will route incoming connections to this port.

The infrastructure has to support this feature and unfortunately Minikube
doesn't. I'll probably add an example later. Until the you can read more about
this in the [official
documentation](https://kubernetes.io/docs/concepts/services-networking/#loadbalancer).

### Using Ingress

I'll look into this in more detail later as an ingress resource has to be set
up first, and I want to learn the basic first.



