# Run the auction server with docker-compose defaults
# Ensures correct ports for PostgreSQL (5434) and Redis (6381)
$env:DB_PORT = "5434"
$env:REDIS_ADDR = "localhost:6381"
$env:PORT = "8081"
go run ./cmd/server
