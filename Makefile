# patient-data-service — atalhos de build/dev.
# No Windows sem `make`, rode os comandos equivalentes à mão (ver README).

MODULE      := github.com/pspd-2026-2-trabalho-2/patient-data-service
PROTO_FILE  := proto/patientdata/v1/patientdata.proto
BIN         := bin/patient-data-service

.PHONY: proto build run test test-it docker compose-up compose-down tidy tools

## Gera o código Go a partir do .proto (requer protoc + plugins no PATH)
proto:
	protoc \
	  --go_out=. --go_opt=module=$(MODULE) \
	  --go-grpc_out=. --go-grpc_opt=module=$(MODULE) \
	  $(PROTO_FILE)

## Instala os plugins protoc do Go
tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

build:
	go build -o $(BIN) ./cmd/server

run:
	go run ./cmd/server

## Testes unitários (sem banco)
test:
	go test ./...

## Testes de integração (precisam do Postgres — ver docker-compose)
test-it:
	go test -tags=integration ./...

docker:
	docker build -t patient-data-service:local .

## Sobe Postgres (schema+seed) + o serviço
compose-up:
	docker compose up -d --build

compose-down:
	docker compose down -v

tidy:
	go mod tidy
