FROM golang:1.18-alpine3.16 AS builder

RUN go version

COPY . /github.com/omekov/superapp-backend/
WORKDIR /github.com/omekov/superapp-backend/

RUN go clean --modcache
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -o apigateway ./cmd/apigateway/*.go

FROM alpine:latest

WORKDIR /root/
COPY --from=0 /github.com/omekov/superapp-backend/apigateway .

CMD ["./apigateway"]
