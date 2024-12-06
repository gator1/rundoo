# Demo: Deploying Apps to Docker Swarm

Swarm is the orchestration component built into the Docker Engine. It's not as widely used as Kubernetes but it's very easy to get started.

## Initialize the Swarm

You don't need any additional installs to use Swarm - just Docker.

```
docker swarm init
```

> You run this command once, and then you can join any number of machines to your Swarm cluster

## Managing the Cluster

In Swarm mode you have additional commands in the Docker CLI.

_The `node` commands let you work with the servers in the cluster:_

```
docker node ls

docker node ls -q

docker node inspect $(docker node ls -q)
```

> Swarms can be multi-architecture: Linux, Windows, Intel and Arm

_The `swarm` commands let you retrieve the token and configure the Swarm:_

```
docker swarm join-token worker

docker swarm join-token manager
```

> In production you would run with 3 or 5 managers - they can support hundreds of workers

_In Swarm mode you run containers as Services, and you group the Services of your app into Stacks:_

```
docker service ls

docker stack ls
```

## Modelling Apps for Docker Swarm

Swarm uses the Docker Compose specification.

You can use a standard [docker-compose.yml](docker-compose.yml) file and run it as a Swarm stack.

_Deploy the WiredBrain app and see what gets created:_

```
docker stack deploy -c docker-compose.yml wb

docker stack ls

docker network ls

docker service ls

docker service ps wb_products-api
```

> On a multi-node Swarm the containers would run across different machines but still connect on the Docker network

Try the app at http://localhost:8080 - it's running in the default development configuration.

## Managing Swarm Apps

You manage stacks and services centrally by connecting to a manager. On a single-node Swarm the manager runs the application workloads too.

_You can manage the application Services through the CLI:_

```
docker service logs wb_web

docker service rm wb_web

docker stack deploy -c docker-compose.yml wb
```

> Swarm uses the same desired-state approach as Docker Compose

But with Swarm the app definition is stored in the cluster.

_Check the details of the web component:_

```
docker service inspect wb_web
```

> You don't need the YAML file to manage the app.

And the stack is the top-level resource, removing it removes all services and networks:

```
docker stack rm wb

docker service ls

docker ps
```