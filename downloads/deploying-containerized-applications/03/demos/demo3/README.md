# Demo: Scale and Reliability in Docker Swarm

Reliability is built-in in Swarm mode - you don't need to add the `restart` option to containers, the Swarm managers make sure your services run at the desired scale.

## Test the Reliability of the App

The app is still running from the previous demo.

_Check the service containers for the stack:_

```
docker stack ls

docker stack ps wb
```

_Verify app failure gets repaired:_

```
docker container ls -f name=wb_stock-api.* -q

docker container exec $(docker container ls -f name=wb_stock-api.* -q) kill 1

docker stack ps wb
```

> Swarm creates a new container - it could run on any node

## Run Services at Scale

There's a new deployment in this [docker-compose.yml](docker-compose.yml) file which enables a debugging feature in the web app.

- the config setting is with an environment variable
- the spec also removes the `depends-on` attributes

_Deploy the update and the web container gets replaced:_

```
docker stack deploy -c docker-compose.yml wb
```

Try the app at http://localhost:8080

> The debug setting shows the host of the web server, which is the ID of the container

You can model Swarm-specific settings in Docker Compose using the `deploy` section - [docker-compose-scale.yml](docker-compose-scale.yml) uses multiple replicas for all the stateless components.

_Scale up the application:_

```
docker stack deploy -c docker-compose-scale.yml wb

docker stack ps wb
```

Try the app at http://localhost:8080 - reload without the local cache (Ctrl-F5 with Vivaldi on Windows), and you'll see different containers respond.

> Swarm uses an ingress network which means multiple containers can listen on the same port - on one node or several nodes.

## Upgrade the Web Component

Container orchestrators provide rolling updates and automated rollbacks.

[docker-compose-v2.yml](docker-compose-v2.yml) upgrades the web component to a new image version.

_Deploy the stack again and print the service containers:_

```
docker stack deploy -c docker-compose-v2.yml wb

docker service ps wb_web 
```

> You can control the speed of the rollout with `update_config` settings in the `deploy` section 

Swarm monitors the upgrade and will take action if the new containers fail. The `v3` web image is misconfigured in the [Dockerfile](../src/web/Dockerfile.v3), so the containers will exit as soon as they start.

_Deploy a new update with the broken image:_


```
docker service update --image psdockerrun/web:v3 wb_web 

docker service ps wb_web 
```

> The default behaviour is to pause a failed rollout - you can configure automatic rollback in your app model

## Teardown

If you're just experimenting with Docker Swarm, you can easily reset everything by forcibly leaving the Swarm. 

That puts you back into standalone Docker Engine mode, and removes all the Swarm resources.

```
docker swarm leave -f
```