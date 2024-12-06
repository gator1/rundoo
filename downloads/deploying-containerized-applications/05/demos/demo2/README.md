# Demo: Deploying apps to ACI

ACI is Azure Container Instances - a managed container service where there is no control plane and no worker nodes. It's a serverless approach to containers.

The WiredBrain app will run in an ACI container group, but we'll define it using Docker Compose and deploy it with the standard Docker CLI.

## Pre-reqs

I'm still using the Postgres database from demo 1. You don't need to create anything in advance to use ACI with Docker - just a resource group.

The Portal [Container Instances page](https://portal.azure.com/#blade/HubsExtension/BrowseResource/resourceType/Microsoft.ContainerInstance/containerGroups) shows nothing right now.


## Connect to an ACI Docker Context

[Docker Contexts](https://docs.docker.com/engine/context/working-with-contexts/) let you control a remote Docker environment from your local CLI.

_Create a context to use ACI:_

```
docker context ls

# docker login azure

docker login azure --tenant-id $tenantId

docker context create aci azure
```

_Switch to use the new ACI context:_

```
docker context ls

# docker context use azure

$env:DOCKER_CONTEXT='azure'
```

_Now Docker commands are running against ACI:_

```
docker ps

docker compose ls
```

### Modelling apps for ACI

You model ACI apps using Docker Compose.

[demo2/docker-compose.yml](docker-compose.yml) is very familiar:

- containers will run in the same container group 
- service location is via DNS using service names
- there's no need for Docker network specs 
- config files can be loaded from an Azure Files share or set using environment variables
- secrets can be uploaded to ACI from local files during the deployment
- port publishing will create a public IP address

### Deploy the app and the secrets

The new context functionality is in the normal `docker` CLI **not** the `docker-compose` CLI.

_Deploy to ACI:_

```
docker compose up

docker compose ls

docker compose ps demo2
```

> Check the container groups in the [Azure Portal](https://portal.azure.com/#blade/HubsExtension/BrowseResource/resourceType/Microsoft.ContainerInstance/containerGroups).

Now the containers are running you can work with them in the usual way.

_Exec into the web container to check secrets:_

```
docker exec -it demo2_web sh

ls /app/secrets

cat /app/secrets/api.json
```

_Check networking from the web container:_

```
wget -O- http://stock-api:8080/stock/1

cat /etc/hosts

exit
```

_Check the application logs:_

```
docker compose logs

docker ps
```

> Browse to web IP address

## Teardown the environment

Removing the app in ACI is as simple as stopping a Docker Compose app locally.

_Remove the ACI container group:_

```
docker compose down

docker ps
```