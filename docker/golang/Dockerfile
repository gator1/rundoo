FROM golang:latest

WORKDIR /usr/src/app

# Install protoc
RUN apt-get update && apt-get install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2


COPY app .
COPY app/go.mod  ./
COPY app/ui ./cmd/portal/ui

RUN go mod download && go mod verify
RUN go install github.com/go-delve/delve/cmd/dlv@latest

#
RUN make compile
RUN go mod tidy
RUN go mod vendor


RUN go build -gcflags="all=-N -l" -v -o /usr/local/bin/app ./cmd/portal...
EXPOSE 2345 
EXPOSE 5001



WORKDIR /usr/src/app/cmd/portal

#CMD ["dlv","debug","--listen=:2345","--headless=true","--api-version=2","--accept-multiclient"]
#CMD ["sh", "-c", "while true; do sleep 3600; done"]
# Use dlv exec to run the pre-built binary in debug mode
CMD ["dlv", "exec", "/usr/local/bin/app", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient"]

