# syntax = docker+earthly:latest

VERSION 0.6

deps:
    FROM golang:1.23.3 
    WORKDIR /go/src/github.com/gator1/rundoo/app

    COPY app/go.mod app/go.sum ./
    RUN go mod download

    COPY app/ .

    # Install protoc and the Go plugins for protobuf and gRPC
    RUN apt-get update && apt-get install -y protobuf-compiler
    RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

build-proto:
    FROM +deps
    RUN protoc --go_out=. --go-grpc_out=. api/v1/rundoo.proto

build-logservice:
    FROM +build-proto
    RUN go build -o /out/logservice ./cmd/logservice

build-portal:
    FROM +build-proto
    RUN go build -o /out/portal ./cmd/portal

build-registryservice:
    FROM +build-proto
    RUN go build -o /out/registryservice ./cmd/registryservice


build-rundooservice:
    FROM +build-proto
    RUN go build -o /out/rundooservice ./cmd/rundooservice

build-api:
    FROM +build-proto
    RUN go build -o /out/api ./api/v1

all:
    FROM +build-logservice
    FROM +build-portal
    FROM +build-registryservice
    FROM +build-rundooservice
    FROM +build-api