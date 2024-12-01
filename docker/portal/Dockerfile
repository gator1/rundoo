
# Tow layers build, first layer is for building the app, second layer is 
# for running the app. can just build the builder like this:
#    docker build --target builder -t my-go-builder .
# and run the app like this:
#    docker run --rm -it my-go-builder
FROM golang:1.23 AS builder
ENV CGO_ENABLED=0

# Install protoc
RUN apt-get update && apt-get install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

WORKDIR /usr/local/go/src/rundoo
COPY app ./app
WORKDIR /usr/local/go/src/rundoo/app
RUN make compile
RUN go mod tidy
RUN go mod vendor

# Ensure the working directory is set correctly for the build command
WORKDIR /usr/local/go/src/rundoo/app
RUN go build -tags=docker -o /portal ./cmd/portal 
# app
FROM alpine:3.13

ENV POSTGRES_CONNECTION_STRING="host=rundoo-db port=5432 user=postgres password=uber dbname=postgres sslmode=disable"
EXPOSE 8080
CMD ["/app/portal"]

WORKDIR /app
COPY --from=builder /portal .


