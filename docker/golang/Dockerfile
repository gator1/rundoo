FROM golang:latest



#COPY go.mod go.* ./
#RUN go mod download && go mod verify
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Install protoc
RUN apt-get update && apt-get install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2




# Download dependencies
#Ensure the working directory is set correctly for the build command
WORKDIR /usr/local/go/src/rundoo
COPY app ./app

# Ensure the working directory is set correctly for the build command
WORKDIR /usr/local/go/src/rundoo/app

ENV GO111MODULE=on


#
RUN make compile
RUN go mod tidy
RUN go mod vendor


RUN go build -gcflags="all=-N -l" -v -tags=docker -o /portal ./cmd/portal 
EXPOSE 2345 
EXPOSE 5001

CMD ["dlv","debug","--listen=:2345","--headless=true","--api-version=2","--accept-multiclient"]




