FROM golang:1.18-alpine3.16 AS builder

RUN go version

COPY . /github.com/omekov/superapp-backend
WORKDIR /github.com/omekov/superapp-backend

RUN go mod download
RUN GOOS=linux go build -o ./.bin/auth ./cmd/auth/*.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/omekov/superapp-backend/.bin/auth .
COPY --from=0 /github.com/omekov/superapp-backend/configs configs/

EXPOSE 80

CMD ["./auth"]
