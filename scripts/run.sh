#!/bin/bash
# Run the auction server with docker-compose defaults
export DB_PORT=5434
export REDIS_ADDR=localhost:6381
export PORT=8081
go run ./cmd/server
