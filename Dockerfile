FROM golang:1.22.5-alpine3.20 AS builder

RUN go version
ENV GOPATH=/

COPY . /github.com/maxzhovtyj/justpay/
WORKDIR /github.com/maxzhovtyj/justpay/

RUN go mod download
RUN GOOS=linux go build -o ./.bin/justpay-linux-amd64 ./cmd/justpay/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 github.com/maxzhovtyj/justpay/.bin/justpay-linux-amd64 .
COPY --from=0 github.com/maxzhovtyj/justpay/configs/config.yml configs/

CMD ["./justpay-linux-amd64", "-config=./configs/config.yml"]
