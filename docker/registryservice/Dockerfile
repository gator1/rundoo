# Tow layers build, first layer is for building the app, second layer is 
# for running the app. can just build the builder like this:
#    docker build --target builder -t my-go-builder .
# and run the app like this:
#    docker run --rm -it my-go-builder
FROM golang:1.20 AS builder
ENV CGO_ENABLED=0

# Install protoc
RUN apt-get update && apt-get install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2


WORKDIR /go/registryservice
COPY ../app/go.mod .
RUN go mod download 

COPY app .
RUN make compile
RUN go build -o /registryservice cmd/registryservice ...
# app
FROM alpine:3.13

ENV POSTGRES_CONNECTION_STRING="host=rundoo-db port=5432 user=postgres password=uber dbname=postgres sslmode=disable"
EXPOSE 8080
CMD ["/app/registryservice"]

WORKDIR /app
COPY --from=builder /registryservice .


