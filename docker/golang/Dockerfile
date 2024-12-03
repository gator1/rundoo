# Use the latest Go image as the builder
FROM golang:latest AS builder

# Install Delve (dlv)
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Install protoc
RUN apt-get update && apt-get install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2


# Enable Go modules
ENV GO111MODULE=on

# Download dependencies
WORKDIR /usr/local/go/src/rundoo
COPY app ./app

# Ensure the working directory is set correctly for the build command
WORKDIR /usr/local/go/src/rundoo/app
RUN make compile
RUN go mod tidy
RUN go mod vendor

#RUN go build -tags=docker -o /portal ./cmd/portal 
# Build the application with debugging flags
RUN CGO_ENABLED=0 GOOS=linux go build -tags=docker  -gcflags="all=-N -l" -o /main ./cmd/portal

# Use the latest Alpine image for the production stage
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built application and Delve from the builder stage
COPY --from=builder /go/bin/dlv /usr/local/bin/dlv
COPY --from=builder /main .
COPY app/ui ./ui

# Expose the necessary ports
EXPOSE 8000 40000

# Set environment variables
ENV PORT=8000

# Command to run the application with Delve
#CMD ["/usr/local/bin/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./main"]
CMD ["sh", "-c", "while true; do sleep 3600; done"]