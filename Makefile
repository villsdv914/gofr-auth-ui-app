# Development commands
run:
	go run cmd/main.go

build:
	go build -o bin/app cmd/main.go

test:
	go test ./...

# Docker commands
docker-build:
	docker build -t gofr-auth-ui-app .

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down

docker-clean:
	docker-compose down -v --remove-orphans

# Database commands
migrate:
	gofr migrate

gen-migration:
	gofr generate migration $(name)

# Utility commands
clean:
	rm -rf bin/
	go clean

deps:
	go mod tidy
	go mod download

# Development setup
setup: deps
	cp .env.example .env 2>/dev/null || echo ".env file already exists"