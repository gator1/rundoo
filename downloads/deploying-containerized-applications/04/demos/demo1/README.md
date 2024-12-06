# Demo: Deploying Apps to Kubernetes

We'll get started with Kubernetes in this demo, deploying the Wired Brain app and seeing how it's modelled.

## Set up the cluster

Setting up Kubernetes is easy if you use a managed service like [AKS](https://azure.microsoft.com/en-gb/services/kubernetes-service/), and much more complex if you do it in the datacenter with [Kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/).

The easiest way to start is with [Docker Desktop](https://www.docker.com/products/docker-desktop) - which gives you a certified single-node Kubernetes cluster.

_Check the node details:_

```
kubectl get nodes

kubectl describe node docker-desktop
```

## Deploy the Wired Brain database

First we'll deploy the database. For that the model is split across two YAML files. The first is:

- [products-db-service.yaml](wiredbrain/products-db-service.yaml) - the network Service

_Create the Service:_

```
kubectl apply -f wiredbrain/products-db-service.yaml

kubectl get service

kubectl get endpoints products-db
```

The next YAML file describes the compute:

- [products-db-deployment.yaml](wiredbrain/products-db-deployment.yaml) - the Deployment

_Create the Deployment and check the Pods:_

```
kubectl apply -f wiredbrain/products-db-deployment.yaml

kubectl get pods

kubectl get endpoints products-db
```

## Deploy the products API

You can group multiple resources in a single YAML file. [products-api.yaml](wiredbrain/products-api.yaml) defines the Service and Deployment for the Java API.

_Create the API resources:_

```
kubectl apply -f wiredbrain/products-api.yaml

kubectl get pods --show-labels

kubectl logs -l app=products-api
```

## Deploy the stock API and website

The remaining components are defined in:

- [stock-api.yaml](wiredbrain/stock-api.yaml) - the Service and Deployment for the Go API
- [web.yaml](wiredbrain/web.yaml) - the Service and Deployment for the .NET website

You can deploy a whole folder of YAML files - Kubernetes takes a desired-state approach and updates or creates what it needs to get to the application model.

_Create the API resources:_

```
kubectl apply -f wiredbrain/

kubectl get svc

kubectl get pods
```

> Now try the app at http://localhost:8080