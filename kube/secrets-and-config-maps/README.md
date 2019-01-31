# Secrets and ConfigMaps

Almost every application needs some kind of configuration. It might be an OAuth
token which is required to authenticate at an external service, username and
password required to connect to a database, or hostnames of other services.
Baking this into a Docker image is possible, but considered bad practice.

Images are usually stored in repositories and on machines which run a Docker
container based on that image. So, if secrets are baked into the images you
can't really control who is using them. Unauthorized people might access the
secrets.

Changing secrets or configurations would require to bake a new image and roll
this out. Also, if there are multiple environments like stage, qa and prod a
Image is deployed to, you would have to bake an image for each of them. 

# When to use ConfigMaps and Secrets

There are several ways to manage secrets and configurations better, and keep
your Docker images portable. the [12 factor app](https://12factor.net/config)
architecture recommends to store them in the environment. Let's see how we can
do this with ConfigMaps and Secrets.

## ConfigMaps

ConfigMaps allow you to store configuration in Kubernetes and provide it to
your service. It supports key-value like configuration, but you can also put
config files into a ConfigMap. Configurations can be provided to a Docker
container as arguments, in the environemnt, or files in a automatically
generated volume.

Let's see how ConfigMaps are defined in yaml:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: simple-http-config
data:
  database.jdbc.url: jdbc://database/users?ssl=true
```

Now, let's create the ConfigMap and take a look at.

```bash
$ kubectl describe cm simple-http-config
Name:         simple-http-config
Namespace:    default
Labels:       <none>
Annotations:  <none>

Data
====
database.jdbc.url:
----
jdbc://database/users?ssl=true
Events:  <none>
```

There is one entry in the data section, which is the key and value we
previously defined in the yaml file.

## Using Environment Variables to Configure an Application

A list of environment variables can be defined for each container in a Pod.

```yaml
kind: Pod
spec:
  containers:
  - image: simple-http:current
    env:
    - name: JDBC_URL
      value: "jdbc://database/users"  
```

The listing above shows how an env variable is defined on the Pod level. With
this approach you need one Pod definition for each enviroment it is supposed to run in.
ConfigMaps decouple a Pod definition from the actual configuration that is
required to run in an enviroment.

```yaml
kind: Pod
spec:
  containers:
  - image: simple-http:current
  env:
  - name: JDBC_URL
    valueFrom:
      configMapKeyRef:
        name: simple-http-config
        key: database.jdbc.url
```

We now use `valueFrom` to set the env variable `JDBC_URL`. The value is taken
from a ConfigMap with name `simple-http-config`. The `key` identifies the entry
of the ConfigMap.

Let's use the yaml in this directory to create both ```kubectl create -f
./simple-http-config.yaml``` and check if the env variable is set to the value
we defined in the ConfigMap. But before we can do this, we have to create a
Service in order to access the Pod. This time with `kubectl`.

```bash
$ kubectl expose pod simple-http-config --type=NodePort --name simple-http-config
service "simple-http-config" exposed

$ minikube service simple-http-config --url
http://192.168.99.111:30574

$ http http://192.168.99.111:30574/env
HTTP/1.1 200 OK
Content-Length: 1152
Content-Type: text/plain; charset=utf-8
Date: Thu, 31 Jan 2019 22:00:53 GMT

{
    "envVars": [
        {
            "key": "HOSTNAME",
            "value": "simple-http-config"
        },
        {
            "key": "JDBC_URL",
            "value": "jdbc://database/users?ssl"
        },
        {
            "key": "SIMPLE_HTTP_CONFIG_PORT_30001_TCP_PORT",
            "value": "30001"
        }
    ]
}
```

And there it is. Great! We are now able to put configuration in ConfigMaps and
use these to set environment variables in our containers. The service can now
be deployed to prod and stage environment and we only need one Docker image and
one Pod definition. But keep in mind that environment variables are only set
once, when a container is started. They don't change when the ConfigMap is
changed. The only way to change the `JDBC_URL` is a deployment, which replaces
the existing Pods with new ones.

You can read more about ConfigMaps in the [official Kubernetes
docs](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#configure-all-key-value-pairs-in-a-configmap-as-container-environment-variables).



