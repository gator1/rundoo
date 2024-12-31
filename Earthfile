VERSION 0.8

all-unit-test:
    BUILD ./app/grpctest/grpc+unit-test
    BUILD ./app/cmd/portal+unit-test
    BUILD ./app/cmd/rundooservice+unit-test

all-docker:
    BUILD ./app/cmd/logservice+docker
    BUILD ./app/cmd/portal+docker
    BUILD ./app/cmd/registryservice+docker
    BUILD ./app/cmd/rundooservice+docker

all-release:
    BUILD ./app/cmd/logservice+release
    BUILD ./app/cmd/portal+release
    BUILD ./app/cmd/registryservice+release
    BUILD ./app/cmd/rundooservice+release

dev-up:
    LOCALLY
    RUN docker-compose up

dev-down:
    LOCALLY
    RUN docker-compose down