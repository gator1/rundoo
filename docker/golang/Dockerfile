#https://stackoverflow.com/questions/63403272/golang-docker-file-for-debug
FROM golang:1.14.7 AS builder
RUN CGO_ENABLED=0 go get -ldflags '-s -w -extldflags -static' github.com/go-delve/delve/cmd/dlv
RUN mkdir /app

ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o main ./...

FROM alpine:3.12.0 AS production
COPY --from=builder /app .
COPY --from=builder /go/bin/dlv /
EXPOSE 8000 40000
ENV PORT=8000
CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./main"]  