.PHONY: build
run:
	go run ./cmd/main.go

build:
	go build -o ./build/main -v ./cmd/main.go

swagger:
	go install github.com/swaggo/swag/cmd/swag@latest && $(HOME)/go/bin/swag init -g ./cmd/main.go --outputTypes json

test:
	go test ./internal/tests/

.DEFAULT_GOAL := run