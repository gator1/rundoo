
# Demo: Running Apps with Docker Compose

Compose is a specification language for describing apps, and a tool to run those apps using Docker.

## Modelling apps in Compose

[docker-compose.yml](./docker-compose.yml) models the same demo app using Compose - it's a clear spec for the component setup.

```
docker-compose up -d
```

> Browse to http://localhost:8080/

Compose uses a desired-state approach, so the commands are repeatable.

```
docker-compose up -d

docker rm -f demo2_products-api_1

docker-compose up -d
```

> Compose will bring the app from the current state to the desired state


## Injecting configuration settings

You can inject configuration with environment variables and volume mounts specified in Compose. 

The [Docker Compose file for test](./test/docker-compose.yml) models the same app in a different environment. Compose lets you run them side-by-side.

```
docker-compose -f ./test/docker-compose.yml up -d

docker network ls

docker ps
```

> The new deployment is at http://localhost:8010 - the original is still at http://localhost:8080

## Modelling reliability

All the configuration you use in `docker` commands is also available in Compose.

[test-v2.yml](./test/test-v2.yml) includes production concerns, limiting compute resources and adding the restart options.

```
docker-compose -f ./test/test-v2.yml up -d
```

> This is an upgrade to the previous deployment - the app is stil at http://localhost:8010/

Containers will restart if they exit when the container process ends.

```
docker exec -it test_stock-api-test_1 kill 1

docker ps --filter name=test*
```

> Try the app again http://localhost:8010

Compose makes the app specification clear and the desired-state approach is easier to work with - but we're still running containers on a single server.
