.PHONY: all build run test clean docker-up docker-down lint

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=server
CMD_PATH=./cmd/server

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(CMD_PATH)

run:
	$(GOCMD) run $(CMD_PATH)

test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

deps:
	$(GOMOD) download
	$(GOMOD) tidy

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

docker-logs:
	docker-compose logs -f

lint:
	golangci-lint run ./...

migrate:
	@echo "Run migrations manually: psql -f migrations/001_initial_schema.sql"
