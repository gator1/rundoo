# Demo: Deploying apps to ECS

ECS is the AWS Elastic Container Service. It's a managed container platform where you can choose to run containers on VMs, or spin them up in a serverless way using AWS Fargate.

You can create a Docker Context to use with ECS and deploy your apps using Docker Compose.

## Pre-reqs

You'll need an AWS account to follow along:

- use the root account or an account with these [IAM roles](https://docs.docker.com/cloud/ecs-integration/#run-an-application-on-ecs)
- use a region with a default VPC

I've already created an [RDS Postgres database](https://us-east-2.console.aws.amazon.com/rds/home?region=us-east-2#databases:), but I don't have any [ECS clusters](https://us-east-2.console.aws.amazon.com/ecs/home?region=us-east-2#/clusters).


## Connect to an ECS Docker Context

_Create an ECS context:_

```
docker context create ecs aws
```

> Choose _AWS environment variables_

Environment variables are a standard way to work with AWS in tools - and they work nicely for automation too.

_Set your AWS environment variables:_

```
$env:AWS_ACCESS_KEY_ID=''
$env:AWS_SECRET_ACCESS_KEY=''
$env:AWS_DEFAULT_REGION='us-east-2'
```

_Switch to the ECS context:_

```
docker context ls

# docker context use aws

$env:DOCKER_CONTEXT='aws'
```

_Check container status:_

```
docker ps

docker compose ls
```

### Secrets & config

You can work with secrets in AWS Secret Manager, deploying them from the Docker CLI.

These local files contain sensitive and standard config settings:

- [api.json](secrets/api.json) - for the web app
- [application.properties](secrets/application.properties) - for the products API

_Create the secrets with the ECS context:_

```
docker secret create wb-products-api ./secrets/application.properties

docker secret create wb-web ./secrets/api.json

docker secret ls
```

> Check on [AWS Secret Manager](https://us-east-2.console.aws.amazon.com/secretsmanager/home?region=us-east-2#!/listSecrets) - you can view and edit secret values


## Deploy the app

The application uses a standard Docker Compose file.

[demo3/docker-compose.yaml](docker-compose.yml) has a few talking points:

- ECS with Docker Context supports service dependency mapping
- environment variables can be used for configuration
- secrets must be mounted from the default path `/run/secrets`
- secrets need to be specified with the AWS Resource Name (ARN)
- port publishing will create a public IP address

> Update the secrets in the Compose file with the real ARNs

_Deploy the app:_

```
docker compose up
```

This will take a minute or two. We'll see progress in the Docker CLI and we can check the [ECS Cluster](https://us-east-2.console.aws.amazon.com/ecs/home?region=us-east-2#/clusters)

The Docker CLI creates a [CloudFormation template](https://docs.docker.com/cloud/ecs-architecture/) and deploys all the pieces to run and publish your app. 

When the app's running you can manage it in the normal way.

_Check the project:_

```
docker compose ls

docker compose ps
```

> Browse to the app

## Teardown the environment

Removing the app in with Compose will remove the ECS cluster and all the associated resources.

_Remove the AWS deployment:_

```
docker compose down
```