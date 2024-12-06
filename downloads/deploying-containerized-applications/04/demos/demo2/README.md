# Demo: Configuring Apps with Kubernetes

Kubernetes has ConfigMap and Secret objects for application configuration. You can model them in YAML or create them from local files.

## Create a new namespace

Namespaces are grouping mechanisms - you might use them to represent different environments or different applications.

Resource names are only unique within a namespace, so you can use the same names for different environments.

[namespace.yaml](wiredbrain/namespace.yaml) defines a new namespace.

_Create a namespace for the test environment:_

```
kubectl apply -f wiredbrain/namespace.yaml

kubectl get namespaces

kubectl get pods --namespace kube-system
```

The default namespace is the default, and the `kube-system` namespace is where all the control plane components run.


## Deploy configuration settings to the test namespace

Resources can specify a namespace in the `metadata` configuration. All the non-sensitive configuration settings for the test environment are modelled in YAML:

- [products-api-properties.yaml](wiredbrain/configMaps/products-api-properties.yaml) - API configuration in Java properties format

- [web-logging.yaml](wiredbrain/configMaps/web-logging.yaml) - web configuration in .NET JSON format

_Create the ConfigMaps in the namespace:_

```
kubectl apply -f wiredbrain/configMaps/

kubectl get configmaps 

kubectl get configmaps -n wb-test
```

ConfigMaps can be viewed to see the plain text.

_Show the data for the app config:_

```
kubectl describe configmap web-logging -n wb-test

kubectl describe cm products-api-properties -n wb-test
```

> The plaintext is visible for anyone with access to the object

## Deploy senstive configuration settings

Secrets are for storing sensitive data. They are not necessarily encrypted in the cluster - depending on the Kubernetes deployment.

The API is similar to ConfigMaps:

- [products-db-password.yaml](wiredbrain/secrets/products-db-password.yaml) stores the database password in a file in a secret

- [stock-api-connection.yaml](wiredbrain/secrets/stock-api-connection.yaml) stores the connection string in a key-value pair

_Create the secrets:_

```
kubectl apply -f wiredbrain/secrets/
```

> You can also create secrets using `kubectl secret create --from-file` if you don't want secrets stored in YAML.

Secrets may be encrypted in the cluster database, but they can be seen from the API server and Kubectl - they're represented in base-64 but not encrypted.

_Check the details of the web app secret:_

```
kubectl -n wb-test describe secret web-api 

kubectl -n wb-test get secret web-api -o jsonpath='{.data.api\.json}'
```

_Decode the base-64:_

```
# add a `base64` function - needed on Windows:
. ./base64.ps1

kubectl -n wb-test get secret web-api -o jsonpath='{.data.api\.json}' | base64 -d
```


## Deploy the app in test configuration

ConfigMaps and Secrets can be loaded into the container filesystem - the Pod references the volume, and the container mounts the volume to the desired path.

- [products-api.yaml](wiredbrain/products-api.yaml) loads secrets and configuration for the Java API

- [web.yaml](wiredbrain/web.yaml) - does the same for the web app

- [products-db.yaml](wiredbrain/products-db.yaml) - loads a secret for the db password

- [stock-api.yaml](wiredbrain/stock-api.yaml) - loads the db connection from a secret into an environment variable


_Deploy the app:_

```
kubectl apply -f wiredbrain/

kubectl get all -n wb-test
```

> Try the app at http://localhost:8081


_The original dev deployment is still running in the `default` namespace:

```
kubectl get po --all-namespaces -l app
```