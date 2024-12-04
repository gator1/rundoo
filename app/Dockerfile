
FROM golang:latest

WORKDIR /usr/src/app

COPY go.mod go.* ./
RUN go mod download && go mod verify
RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY . .
RUN go build -a -gcflags="all=-N -l" -v -tags="docker" -o  /usr/local/bin/app ./cmd/portal 
EXPOSE 2345 
EXPOSE 5001

WORKDIR /usr/src/app/cmd/portal

CMD ["dlv","debug","--listen=:2345","--headless=true","--api-version=2","--accept-multiclient"]