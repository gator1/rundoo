
# Demo: Modelling Scale with Docker Compose

Compose is valid for some production scenarios, but it is limited.

## Scaling with the CLI

You can scale up and run multiple containers for a component - but only if the component doesn't use restricted resources like ports.

[This Docker Compose file](./docker-compose.yml) publishes the stock API to a random host port.

```
docker-compose up -d

docker-compose up -d --scale stock-api-test=3
```

> Browse to http://localhost:8010, requests from the web to the Stock API are load-balanced

```
docker-compose logs stock-api-test
```

But you can't scale containers that publish ports, because ports are an exclusive resource.

```
docker-compose up -d --scale web-test=3

docker ps -a
```

> New web containers are created but can't start; stock API scaled down

This is the other problem with Compose - it's a client-side component. The server doesn't store the application model, and you need to have the YAML available to manage the app.

## Recording scale in the model 

The [v2.yml](./v2.yml) specifies the API scale, so it's part of the model and you won't accidentally modify the current state.

`scale` is an older part of the Compose specification - the latest releases of the Compose CLI let you remove the version tag in your YAML.

```
docker-compose -f ./v2.yml up -d
```

> The failed web containers are removed, and new Stock API containers are created

But the real problem with the Compose CLI is that it's not a container platform. It works for a single server - there's no way to cluster it and that means no shared data store. Without a clustered data store you don't get features like stored configuration and secrets.