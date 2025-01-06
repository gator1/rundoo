VERSION 0.8

all-unit-test:
    BUILD ./app/grpctest/grpc+unit-test
    BUILD ./app/cmd/portal+unit-test
    BUILD ./app/cmd/rundooservice+unit-test

all-docker:
    BUILD ./app/+docker-logservice
    BUILD ./app/+docker-portal
    BUILD ./app/+docker-registryservice
    BUILD ./app/+docker-rundooservice

all-release:
    BUILD ./app/+build-logservice
    BUILD ./app/+build-portal
    BUILD ./app/+build-registryservice
    BUILD ./app/+build-rundooservice

dev-up:
    LOCALLY
    RUN docker compose up

dev-down:
    LOCALLY
    RUN docker compose down