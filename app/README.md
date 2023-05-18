# Rundoo: A Fully Working Distributed System for Rundoo Products

![Go Logo](https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Aqua.png)

## Rundoo

A fully working distributed system for Rundoo products. This system consists of a registry service, a log service, rundoo service (business logic), and a web portal.

- Website: 
- Instruction: [https://you.ashbyhq.com/rundoo/assignment/48c21a42-c420-48bb-8f33-45cf12bec30d](https://you.ashbyhq.com/rundoo/assignment/48c21a42-c420-48bb-8f33-45cf12bec30d)
- Email: [ilovegators@gmail.com](mailto:ilovegators@gmail.com)

## Overview

Rundoo is a distributed system developed in Go. The codebase does not rely on JavaScript, even for the front-end. The service utilizes Go's built-in HTTP server. For a real-world application, it is recommended to consider frameworks such as Gin or Uber Fiber.

The communication between the web portal and rundoo service is done through gRPC, while the other services communicate using REST APIs over HTTP. Initially, the plan was to use UberFx and gRPC for all communication, but due to closed-source modules in UberFx, the development time would have been significantly longer. Nevertheless, this project showcases expertise in a distributed system environment using pure Go.

## UI Overview

The web portal is accessible at `localhost:5050` and provides functionality to view, add, and search products. When the database is initially empty, you can use the provided SQL files to populate it. The portal consists of three buttons: "Home," "Add Product," and "Search Product." The "Home" button allows you to view the products in the database, while the "Add Product" and "Search Product" buttons enable you to add new products or search for existing ones. Please note that the current implementation lacks pagination, making it suitable for demo environments rather than large-scale production use.

## Database Setup

The rundoo service utilizes a PostgreSQL database, which can be set up using the following steps:

1. Pull the PostgreSQL Docker image:
   ```sh
   docker pull postgres


This gives you the Postgres image. I have 
```sh
 docker images
REPOSITORY   TAG       IMAGE ID       CREATED      SIZE
postgres     latest    fea618ce20fd   6 days ago   360MB
```

2. Run Postgres in a container

```sh
docker run --name rundoo-db-container -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 postgres 
```
This sets up the admin password for Postgres. It also builds the port forwarding so you can access the database. 

3. Install the command line tool, psql
```sh
brew install postgresql
```

4. Logon to the database
```sh
psql -h localhost -p 5432 -U postgres
```

5. Setup database 
```sh
CREATE DATABASE rundoo; 
CREATE ROLE rundoo WITH LOGIN PASSWORD 'uber';
\c rundoo;
```

6. Setup table
```sh
CREATE TABLE IF NOT EXISTS products (
id bigserial PRIMARY Key,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
name text NOT NULL,
category text NOT NULL,
sku text NOT NULL,
version integer NOT NULL DEFAULT 1
);
```

There is a sql script in the codebase called setup.sql. It doesn't seem to work to create the database. You can try. I think rest of the commands work well. 
```sh
cd ~/HOME/your path/app
psql -h localhost -p 5432  -U postgres -d rundoo -f setup.sql
```
You can also fill in some entries for the database with mock datain mockdata.sql. 
```sh
psql -h localhost -p 5432  -U postgres -d rundoo -f mockdata.sql
```


Additionally, there is a `setup.sql` script in the codebase that you can try using to create the database and run the other commands mentioned above. You can also populate the database with mock data using the `mockdata.sql` script.

## Database Optimization

To speed up database queries, indexing is crucial. In this case, you can create indexes on the frequently searched columns. Connect to the database and execute the following SQL commands to create the indexes:

```sh
CREATE INDEX idx_name ON products (name);
CREATE INDEX idx_category ON products (category);
CREATE INDEX idx_sku ON products (sku);
CREATE INDEX idx_name_category_sku ON products (name, category, sku);
```
When the database grows larger, optimizing queries with appropriate indexes will significantly improve performance.
Sharding could also be used, we could use products alphabet to build sharding. 
We could also use the version to sunset old entries to a datalake. In the current table design I didn't have a timestamp. We could have one for the production system. We could use the timestamp to shard the table and index. For example we could do the sharding by month. 

## Service operation
Go code is meant to be built. But during my development I find it to be easier to use go run to test. In real life the deployment should be container based but here I am running them in a terminal. In Vscode it's easy to use cmd+\ to split the terminal window. I have four terminals to run the services in the order listed below. You need to be at the app directory. Before running go build, you need to run 
```sh
make compile
```
in the app directory to generate the grpc files. 

```sh
go run cmd/registryservice/main.go
go run cmd/logservice/main.go
go run $(ls cmd/rundooservice/*.go | grep -v _test.go)
go run $(ls cmd/portal/*go | grep -v _test.go)
```

From the window that runs registry it will print out if a service is up or down; this is what a registry service does in real life. You will also see the heartbeat messages.

The log service makes all the logs from all the services to go to it and save them in app/app.log. It's the first place to look if you run into any issue. This is similar to a production environment; the logs could be processed for ELK. 


## The design 
Registry service receives each service registration, a service endpoints are recorded and can be quried by other services. A heartbeat is sent periodically.   

There are four services:
Registry: service registration and health monitoring. 
Log service: Centralized logging. 
Rundoo Service: Business logic and data persistence
Portal: Web application and API gateway. 

## Testing
Unit tests and integration tests are provided in the codebase. You can run the tests using the following command:

```sh
go test ./...
```

## Deployment (future)
To deploy the Rundoo distributed system, you can follow these general steps:

Set up the required infrastructure (servers, load balancers, etc.) based on the system's architecture.

Clone the Rundoo repository:

```sh
your repo
```
Build the necessary binaries:

```sh
go build ./...
```

Configure the system by updating the configuration files according to your environment.

Deploy and start the individual services in the appropriate order:

Registry service
Log service
Rundoo service
Web portal
Ensure the services are running and accessible.

## Limitation

Unit tests are minimum. In real life more unit tests should be done to make sure all the code paths are covered. 



## FAQ

