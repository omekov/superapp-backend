FROM golang:1.23-alpine3.20 AS builder

RUN go version

COPY . /github.com/omekov/dubcaicar/
WORKDIR /github.com/omekov/dubcaicar/

RUN go clean --modcache
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dubaicarkz ./cmd/dubaicarkz/*.go

FROM alpine:latest

WORKDIR /root/
COPY --from=0 /github.com/omekov/dubcaicar/dubcaicar .

CMD ["./dubcaicar"]