# ReplicaSets

A ReplicaSet ensures that a specified number of pod replicas are running. It
replaces the Replication Controller, which is deprecated and will eventually
removed.

When a pod fails, isn't reachable, or an entire Node fails, a ReplicaSet
recreates new Pods automatically. Additionally, a liveness check can be defined,
which is used by a ReplicaSet to check if a service is still alive. If a JDK process
is `OutOfMemory` for example, the process and thus the container are still running,
but is will not responde to HTTP requests anymore. A HTTP GET liveness check helps to
identify affected pods. A ReplicaSet deletes the pod and creates a new one.

The number of replicas can be increased and decreased manually. ReplicaSets
support [Horizonal
PodAutoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
and can scale a pod up and down based CPU utilization.

## Creating ReplicaSets

Let's take a look at a simple YAML for a ReplicaSet.

```yaml
apiVersion: apps/v1beta2
kind: ReplicaSet
metadata:
  name: simple-http-rs
spec:
  replicas: 2
  selector:
    matchLabels:
      app: simple-http-rs
  template:
    metadata:
      labels:
        app: simple-http-rs 
    spec:
      containers:
      - image: simple-http:current
        name: simple-http
        ports:
        - containerPort: 30000
          protocol: TCP
        livenessProbe:
          httpGet:
            path: /health
            port: 30000
          initialDelaySeconds: 15
```

The ReplicaSet ensures that 2 `replicas` of a pod are running at any time.

The `selector` is used to identify the pods overseen by a ReplicaSets. And this
is a main improvement of a ReplicaSet over ReplicationControllers, because it
supports [set-based
selectors](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#resources-that-support-set-based-requirements),
which are more powerful than the equality-based selectors (`matchLabel`).
They allow to define several expressions which are used to idenfity managed pods.
 
```yaml
matchExpressions:
  - key: app
    operator: In
    values:
    - simple-http-rs
    - simple-http-rs-2
```

The `matchExpressions` in the yaml above would replace the `matchLabels` in the
ReplicaSet spec. In this example, all pods having either `simple-http-rs` or
`simple-http-rs-2` would be matched.

The `template` contains a pod specification, which is used by a ReplicaSet to
create new pods.

`livenessProbe` defines a HTTP GET health check for the `simple-http` service.

## ReplicaSets in Action

Let's create a ReplicaSet:

```bash
$ kubectl create -f ./simple-http.yaml
replicaset.apps "simple-http-rs" created
```

Now let's check if the specified pod replicas are running:

```bash
$ kubectl get pods
NAME                     READY     STATUS    RESTARTS   AGE
simple-http-rs-chtvt     1/1       Running   0          1m
simple-http-rs-ww2r4     1/1       Running   0          1m
```

Two `simple-http` pods are running. Let's now look how to sent a HTTP request
to the simple-http service running in a pod.

```bash
$ kubectl port-forward simple-http-rs-ww2r4 30000:30000
$ http :30000/health

HTTP/1.1 200 OK
Content-Length: 50
Content-Type: text/plain; charset=utf-8
Date: Sat, 26 Jan 2019 20:49:23 GMT

{
    "hostname": "simple-http-rs-ww2r4",
    "status": "OK"
}

```

Note that the first command forwards the local port 30000 to port 30000 of a
specific pod. If this pod now disappears, requests would run into connection
errors.. ReplicaSet would create a new pod, but this specific pod would never
come back. So, in real applications you should never directly connect to a pod
unless low availability is ok. Use Kubernetes Services instead to connect to
pods.

Now let's delete one pod. The ReplicaSet should immediately create a new pod.

```bash
$ kubectl delete pod simple-http-rs-chtvt
pod "simple-http-rs-chtvt" deleted

$ kubectl get pods
NAME                     READY     STATUS        RESTARTS   AGE
simple-http-rs-chtvt     0/1       Terminating   0          2m
simple-http-rs-ww2r4     1/1       Running       0          2m
simple-http-rs-z8wlv     1/1       Running       0          2s
```

As expected, the ReplicaSet created a new pod, which is already running.

## When to use ReplicaSets

You typically don't create ReplicaSets on your own to deploy pods.
Instead you would typically use a Deployment, which internally creates a
ReplicaSet.

ReplicaSets don't support the `rolling-update` command. This is something you
definitly want to have for a highly available application. Even if it would
support this imperative command, Deployments are declarative and thus
recommended to use.

If the process a Deplyoment use to deploy your application doesn't work for
you, a ReplicaSet can be used to build something which better fits your needs.


