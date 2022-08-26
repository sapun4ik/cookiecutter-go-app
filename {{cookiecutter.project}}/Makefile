.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/main.go

run:
	go run cmd/main.go 

test:
	go test --short -coverprofile=cover.out -v ./...
	make test.coverage

test.coverage:
	go tool cover -func=cover.out | grep "total"

swag:
	swag init -g internal/app/app.go

lint:
	golangci-lint run
# golangci-lint run --new-from-rev master

lint-fix:
	golangci-lint run --fix


# ==============================================================================
# Docker

local:
	@echo Clearing prometheus data
	rm -rf ./prometheus
	@echo Starting local docker compose
	docker-compose -f docker-compose.local.yaml up --build $(demon)

locald: 
	make demon=-d local 

develop:
	@echo Clearing prometheus data
	rm -rf ./.tmp/prometheus
	@echo Starting local docker compose
	docker-compose -f docker-compose.yaml up --build $(demon)

developd:
	make demon=-d develop 

# ==============================================================================
# Docker support

FILES := $(shell docker ps -aq)

down:
	docker stop $(FILES)
	docker rm $(FILES)

clean:
	docker system prune -f

logs:
	docker logs -f $(FILES)