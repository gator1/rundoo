# Demo: Configuring Apps with Docker Swarm

Swarm managers run a distrbuted database which stores the cluster state and application models. You can use the data store for application configuration.

## Create Config Objects

We'll run the app in the test configuration. For that we need to deploy it with new config files:

- [logging](configs/web/logging.json) for the .NET web app
- [app settings](configs/products-api/application.properties) for the Java products API

Anything non-senstive can be stored in a Swarm Config object.

_Create the application Configs:_

```
docker config create products-api-config configs/products-api/application.properties

docker config create web-logging configs/web/logging.json
```

> Config objects can be any file format, and the contents are available to view

_Print the details and the content of the Configs:_

```
docker config inspect --pretty products-api-config 

docker config inspect --pretty web-logging
```

## Create Secret Objects

Sensitive data can be stored in Secrets - the contents aren't visible through the CLI or the Docker API.

We'll use secrets for:

- [API settings](secrets/web/api.json) for the .NET web app
- [database connection](secrets/products-api/db.properties) for the Java products API
- [password](secrets/products-db/POSTGRES_PASSWORD) for the Postgres database

_Create the Secrets:_

```
docker secret create products-api-dbconfig secrets/products-api/db.properties

docker secret create web-api secrets/web/api.json

docker secret create products-db-password secrets/products-db/POSTGRES_PASSWORD
```

Secrets are encrypted inside the database and in transit through the cluster. They're only surfaced as plain text inside containers which are configured to use them.

_Print the details of the database Secret:_

```
docker secret inspect --pretty products-db-password
```

## Modelling Configuration in Docker Swarm

The new [docker-compose.yml](docker-compose.yml) file includes `configs` and `secrets`:

- they're loaded into services at a specific file path
- and they're defined as external, so they need to exist before deployment

_Check all the dependencies are there:_

```
docker config ls

docker secret ls
```

_Now deploy the stack:_

```
docker stack deploy -c docker-compose.yml wb
```

> This workflow means separate teams can own app deployment and configuration management

Check the app at http://localhost:8080 - this is the test configuration.

## See How Secrets and Configs are Loaded

You can connect to a container from the node that's running it. There you can inspect the filesystem to see how configuration gets loaded.

_Exec into a shell session for the Java API:_

```
docker container ls -f name=wb_products-api.*

docker container ls -f name=wb_products-api.* -q

docker container exec -it $(docker container ls -f name=wb_products-api.* -q) sh
```

_And print the contents of the config files:_

```
ls /app

ls -l /app/config

cat /app/config/application.properties

cat /app/config/db/application.properties
```

> Inside the container you see the plain-text file contents