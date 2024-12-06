# Demo: Deploying Kubernetes apps to AKS

AKS is the Azure Kubernetes Service - a managed Kubernetes service where you pay for the worker nodes but not the control plane.

We'll deploy the same WiredBrain app from the Kubernetes model, but instead of running a database container, we'll connect the app to a managed Azure Postgres database.

## Pre-reqs

For this demo the cluster and database have already been created using these scripts:

- [aks.ps1](./azure/aks.ps1) - to create the cluster
- [postgres.ps1](./azure/postgres.ps1) - for the database
- [init-products-db.sql](./database/init-products-db.sql) - to deploy the schema & reference data

You can find those resources in the Azure Portal - [Azure Database for PostgreSQL servers
](https://portal.azure.com/#blade/HubsExtension/BrowseResource/resourceType/Microsoft.DBforPostgreSQL/servers) and [Azure Kubernetes Service](https://portal.azure.com/#blade/HubsExtension/BrowseResource/resourceType/Microsoft.DBforPostgreSQL/servers).


> To follow along you'll need an Azure subscription and the `az` command line, and a SQL client like [Sqlectron](https://sqlectron.github.io).

## Connect to the cluster

Kubectl uses a configuration file to connect to a cluster, you can switch between clusters using **contexts**.

_List current contexts:_

```
kubectl config get-contexts
```

_Add the details for an AKS cluster to your Kubectl config:_

```
az aks list -o table 

az aks get-credentials -g ps-docker-run -n ps-docker-run-aks

kubectl config get-contexts
```

_Swtich to the AKS cluster:_

```
kubectl config use-context ps-docker-run-aks

kubectl get nodes
```

## Modelling apps for AKS

The application model in the YAML is pretty much the same. These are the differences:

- there's no Service or Deployment for the database - that's an external service now
- there's no YAML for the Secrets, we'll use a different workflow to deploy them
- the Deployments have some production concerns modelled as in [products-api.yaml](wiredbrain/products-api.yaml)

## Create the Secrets

Secrets need to exist before the app is deployed, because the Deployments specify them.

- [create-secrets.ps1](secrets/create-secrets.ps1) does that, using local files
- [application.properties](secrets/sources/application.properties) contains the database setup for the Java API

_Create the namespace and secrets:_

```
kubectl apply -f wiredbrain/namespace.yaml

cd secrets

./create-secrets.ps1

kubectl get secrets -n wb-prod
```

## Deploy the app

The rest of the deployment is standard Kubernetes.

_Deploy the configuration and the app:_

```
kubectl apply -f wiredbrain/configMaps/

kubectl apply -f wiredbrain/
```

The Service for the web app is a LoadBalancer - AKS creates a public IP address for the service.

_Check the public IP:_

```
kubectl get svc -n wb-prod
```

> Browse to the external IP address