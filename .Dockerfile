FROM golang:1.23rc1-alpine3.19 AS builder

RUN go version

COPY . /github.com/omekov/dubcaicar/
WORKDIR /github.com/omekov/dubcaicar/

RUN go clean --modcache
RUN  go mod download && \
     CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -o dubaicarkz ./cmd/dubaicarkz/*.go

FROM alpine:latest

WORKDIR /root/
COPY --from=0 /github.com/omekov/dubcaicar/dubcaicar .

CMD ["./dubcaicar"]