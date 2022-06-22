# See README.txt.

.PHONY: run init proto clean docker

run:
	go run ./cmd/auth/*.go

test:
	go test -v -cover -timeout 60s ./...

init:
	-mkdir -p docs
	-rm -rf ./vendor
	-go install github.com/pressly/goose/v3/cmd/goose@latest
	-go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	-go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	-go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	-go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	-go install github.com/favadi/protoc-go-inject-tag@latest
	-go install github.com/golang/mock/mockgen@latest
	-go install github.com/matryer/moq@latest
	-go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go mod tidy
	go mod vendor
	make generate
	go mod vendor

generate: protobuf
	find $(CWD)/pkg/proto -type f -name "mock.go" -delete
	go generate ./...

proto:
	protoc --proto_path=proto proto/auth/v1/*.proto --go_out=internal/auth/delivery/grpc/v1 --go-grpc_out=internal/auth/delivery/grpc/v1

clean:
	rm -rf ./internal/auth/v1/proto/*.go

mocksrepository:
	mockgen -source=./internal/auth/user/repository/repository.go -destination=./internal/auth/user/repository/mocks/repository.go -package=mocksrepository

docker:
	make local-db 
	make local-rdb
	make goose-auth

local-db:
	-docker kill superapp_postgres
	-docker run -d --rm --name=superapp_postgres \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_USER=superapp_user \
		-e POSTGRES_DB=superapp_auth \
		-p 3432:5432 \
		postgres:13-alpine
	

local-rdb:
	-docker kill superapp_redis
	-docker run -d --rm --name superapp_redis \
		-e REDIS_PASSWORD=superapp \
		-p 6379:6379 redis /bin/sh -c 'redis-server --appendonly yes --requirepass superapp'

lint:
	golangci-lint run

goose-auth:
	goose -dir ./migrations/auth postgres "user=superapp_user password=password dbname=superapp_auth sslmode=disable port=3432" up
	make goose-auth-directory

goose-auth-directory:
	goose -dir ./migrations/auth/directory postgres "user=superapp_user password=password dbname=superapp_auth sslmode=disable port=3432" up

