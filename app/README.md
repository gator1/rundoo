A fully working distributed system for Rundoo products.
============================

- Website:https://github.com/gator1/rundoo
- Instruction:https://you.ashbyhq.com/rundoo/assignment/48c21a42-c420-48bb-8f33-45cf12bec30d
- email: [Gmail](hilovegators@gmail.com)

<img src="https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Aqua.png" width="600px">

# Rundoo

It's a distributed system that a registry service, a log service, rundoo service(business logic), and a web portal. 

The code is developed in Go; there is no Javascript even for front end. 
The service uses Go httpserver. For a real application one probabaly should consider
some framework such as Gin or Uber fiber. 

The web portal and rundoo service commuinate through grpc while all other services talk with REST API over http. 
I planed to use UberFx to use gRPC for all communications. However I realized that Uber didn't open source some of modules
so using UberFx as a foundation would take longer time to build; the system would be much better but I just didn't have enough time. 
So this is a pure Go play that hopefully demonstrate my experitises in a distributed system envrionment. 


## UI overview

The Web portal listens on port 5050 (on Mac 5000 is often taken). It will show the products in the DB; when you start the DB is empty. Some sql files are provided to help you get some data in. 

There are three buttons, Home, Add Product, and Search Product. 

## Database

The rundoo service uses a Postgres database. The Postgres runs in a docker.  The instruction to set up Postgres, it's assumed that you have docker installed. All my tests are done on a Mac. 

1. Pull Postgres imagine
```sh
docker pull postgres
```
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
This set up the admin password for Postgres. It also builds the forwarding so you can access the database. 

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
You can also fill in some entries for the database
```sh
psql -h localhost -p 5432  -U postgres -d rundoo -f mockdata.sql
```


## Database speed up

Indexing is the first thing. I login to the database and nbuild these indexing. 
```sh
CREATE INDEX idx_name ON products (name);
CREATE INDEX idx_category ON products  (category);
CREATE INDEX idx_sku ON products (sku);
CREATE INDEX idx_name_category_sku ON products (name, category, sku);
```

When the database becomes big, sharding could be used, we could use products alphabet to build sharding. 
We could also use the version to sunset old entries to a datalake. In the current table design I didn't have a timestamp. We could have one for the production. We could use the timestamp to shard the table and index. For example we could do the sharding by month. 

## Service operation
Go code is meant to be built. But during my development I find it easier to use go run to test. In real life the deployment should be container based but here I am running them. In Vscode it's easy to use cmd+\ to split the terminal. I have four terminals to run the service in the order. You need to be at the app directory. Before running go build, you need to run 
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

From the window that runs registry it will print out if a service is up or down; this is what a registry service does in real life. 

The log service makes all the logs from all the services to go to it and save them in app/app.log. It's the first place to look if you run into any issue. This is similar to a production environment; the logs could be processed for ELK. 


## Behind the scenes 
Registry service receives each service registration, a service endpoints are recorded and can be quried by other services. A heartbeat is sent periodicly  


## Limitation

NUnit tests are minimum. In real life more unit tests should be done to make sure all the code paths are covered. 
, in real life there should be enough unit tests to cover all code path. 



## FAQ

