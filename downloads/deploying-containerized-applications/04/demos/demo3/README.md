# Demo: Scale and Reliability in Kubernetes

Kubernetes provides reliability at different levels, using the compute abstractions - Pods, Deployments and ReplicaSets.

The same mechanisms support running at scale and safely deploying application upgrades.

## Deploy the Wired Brain app at scale

We'll use a dedicated namespace for this demo - the other deployments can still be running without affecting this one.

_Deploy the new test configuration settings:_

```
kubectl apply -f wiredbrain/namespace.yaml

kubectl apply -f wiredbrain/configMaps/

kubectl apply -f wiredbrain/secrets/

kubectl get cm,secrets -n wb-test-2
```

This deployment will run at scale using the `replicas` value in the Deployment spec.

It also sets up the debug configuration to show the hostname of the web server - that's modelled as an environment variable in [web.yaml](wiredbrain/web.yaml).

_Deploy the remaining app resources:_

```
kubectl apply -f wiredbrain/

kubectl get pods -n wb-test-2 --watch

kubectl get pods,svc -n wb-test-2

```

> Browse to http://localhost:8082, refresh to see load-balancing in action


## Verify the resilience of the app

Pods are responsible for keeping their containers running. If a container exits, the Pod restarts it.

_Simulate an app failure in the stock API:_

```
kubectl get deploy stock-api -n wb-test-2

kubectl get po -l app=stock-api -n wb-test-2

kubectl exec deploy/stock-api -n wb-test-2 -- kill 1

kubectl get po -l app=stock-api -n wb-test-2
```

> One Pod shows a restart. The container doesn't get replaced - the same container gets restarted.


## Perform an automated update of the website

[update-v2/web.yaml](wiredbrain/update-v2/web.yaml) updates the web Deployment to use a new Docker image.

The Deployment creates a new ReplicaSet with the updated configuration, and incrementally manages the scale - scaling the new ReplicaSet up and the old ReplicaSet down.

_Deploy the update:_

```
kubectl apply -f wiredbrain/update-v2/; `

kubectl get replicasets -n wb-test-2 -l app=web --watch
```

> At the end of the rollout the old RS has a scale of 0 and the new one a scale of 3

Try the app again at http://localhost:8082, the new UI comes from the new Pods.


## Verify that failed updates don't take the app offline

Services are loosely coupled to Pods using a label selector. Any Pods that match the label are candidates for Service endpoints, but the Service only sends traffic to healthy Pods.

That means if your application update fails, the new Pods won't get any traffic and the old ReplicaSet won't be scaled down.

_Deploy an update to the failing v3 web image:_

```
kubectl apply -f wiredbrain/update-v3/ ;`

kubectl get rs -l app=web -n wb-test-2 --watch
```

> While the update is in progress, the v2 Pods keep responding at http://localhost:8082

_Check the status of the Pods:_

```
kubectl get pods -l app=web,version=v3 -n wb-test-2

kubectl describe pods -l app=web,version=v3 -n wb-test-2
```

Kubernetes keeps restarting the Pod, hoping it will work. There's a backoff period between restarts so they become spaced out over minutes to protect the node.
