.PHONY: build run test docker-up docker-down

build:
	go build -o bin/rate-limiter

run:
	go run main.go

test:
	go test ./... -v

test-short:
	go test ./... -v -short

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

docker-test:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml up --build test

.DEFAULT_GOAL := run