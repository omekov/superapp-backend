# See README.txt.

.PHONY: run-auth run-api test init proto clean docker cover proto-gw

run-auth:
	go run ./cmd/auth/*.go

run-api:
	go run ./cmd/apigateway/*.go

test:
	go test -v -cover -timeout 60s ./...

cover:
	go test -short -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

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
	-go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go mod tidy
	go mod vendor
	make generate

generate: protobuf
	find $(CWD)/pkg/proto -type f -name "mock.go" -delete
	go generate ./...

proto-auth:
	rm -f internal/auth/delivery/grpc/v1/*.go
	protoc --proto_path=proto/auth --go_out=internal/auth/delivery/grpc --go_opt=paths=source_relative \
	--go-grpc_out=internal/auth/delivery/grpc  --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=internal/auth/delivery/grpc --grpc-gateway_opt=paths=source_relative \
	proto/auth/v1/*.proto

proto-api:
	rm -f internal/apigateway/delivery/grpc/v1/*.go
	protoc --proto_path=proto/apigateway --go_out=internal/apigateway/delivery/grpc --go_opt=paths=source_relative \
	--go-grpc_out=internal/apigateway/delivery/grpc  --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=internal/apigateway/delivery/grpc --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=apigateway \
	proto/apigateway/v1/*.proto 
	statik -src=./doc/swagger -dest=./doc

clean:
	rm -rf ./internal/auth/v1/proto/*.go

mocks:
	- mocks-auth-user-repository
	- mocks-auth-user-service

mocks-auth-user-repository:
	mockgen -source=./internal/auth/user/repository/repository.go -destination=./internal/auth/user/repository/mocks/repository.go -package=mocksrepository

mocks-auth-user-service:
	mockgen -source=./internal/auth/user/service/service.go -destination=./internal/auth/user/service/mocks/service.go -package=mocksservice

docker:
	make local-db 
	make local-rdb

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

