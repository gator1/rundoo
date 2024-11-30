FROM golang:1.15.6 AS builder
ENV CGO_ENABLED=0

WORKDIR /go/registryservice
COPY ../app/go.mod .
RUN go mod download 

COPY ../app .
RUN make compile
RUN go run cmd/registryservice/main.go

# app
FROM alpine:3.13

ENV POSTGRES_CONNECTION_STRING="host=rundoo-db port=5432 user=postgres password=uber dbname=postgres sslmode=disable"
EXPOSE 8080

