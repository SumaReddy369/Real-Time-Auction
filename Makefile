.PHONY: build run test docker-up docker-down integration-test

build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

test:
	go test -v ./...

docker-up:
	docker-compose up -d
	@echo "Waiting for services..."
	@sleep 5

docker-down:
	docker-compose down

integration-test: docker-up
	@echo "Running integration tests..."
	@powershell -ExecutionPolicy Bypass -File ./scripts/integration_test.ps1
	$(MAKE) docker-down
