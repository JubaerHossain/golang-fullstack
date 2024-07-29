APP_NAME = "go-echo-template"
install:
	go mod tidy
	bun install
	bun run dev:css

dev:
	go run cmd/main.go

build:
	go build -o bin/$(APP_NAME) cmd/main.go

run:
	/bin/bash -c "bin/$(APP_NAME)"