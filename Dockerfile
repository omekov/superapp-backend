FROM golang:1.18-alpine3.16 AS builder

RUN go version

WORKDIR /github.com/omekov/superapp-backend/
COPY . .

RUN go clean --modcache
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -o auth ./cmd/auth/*.go

FROM alpine:latest

COPY --from=0 /github.com/omekov/superapp-backend/auth /

CMD ["./auth"]
