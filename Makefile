.PHONY: build test lint docker-build docker-up

BIN_DIR := bin
BIN := $(BIN_DIR)/blobtube

build:
	@mkdir -p $(BIN_DIR)
	go build -trimpath -o $(BIN) ./cmd/server

test:
	go test ./...

lint:
	golangci-lint run

docker-build:
	docker build -t blobtube:dev .

docker-up:
	docker compose up --build
