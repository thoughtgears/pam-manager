.PHONY: dev lint

lint:
	golangci-lint run

dev:
	godotenv -f .env go run main.go
