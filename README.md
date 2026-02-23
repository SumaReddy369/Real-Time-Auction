# Real-Time Auction System

A distributed real-time auction system built with Go, featuring WebSockets for live bid updates, Redis Pub/Sub for broadcasting, and PostgreSQL for persistence.

## Features

- **REST API** for auction and bid management
- **WebSocket** for real-time bid updates
- **Redis Pub/Sub** for broadcasting bids across server instances
- **PostgreSQL** for durable storage
- **Background worker** for automatic auction expiry

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auctions` | Create a new auction |
| GET | `/api/v1/auctions/:id` | Get auction details |
| POST | `/api/v1/auctions/:id/bids` | Place a bid |
| GET | `/api/v1/auctions/:id/bids` | List all bids for an auction |
| WS | `/ws/auctions/:id` | WebSocket for real-time bid updates |

## Tech Stack

- **Go 1.21+**
- **Gorilla WebSocket** - WebSocket support
- **Gorilla Mux** - HTTP router
- **Redis** - Pub/Sub for real-time broadcasting
- **PostgreSQL** - Auction and bid storage
- **pgx** - PostgreSQL driver

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose (for PostgreSQL and Redis)

### 1. Start dependencies

```bash
docker-compose up -d
```

### 2. Run the server

With docker-compose (uses ports 5434, 6381, 8081):

```bash
# Windows
.\scripts\run.ps1

# Linux/Mac
./scripts/run.sh
```

Or with env vars:

```bash
DB_PORT=5434 REDIS_ADDR=localhost:6381 PORT=8081 go run ./cmd/server
```

Or build and run:

```bash
go build -o bin/server ./cmd/server
DB_PORT=5434 REDIS_ADDR=localhost:6381 PORT=8081 ./bin/server
```

### 3. Create an auction

```bash
curl -X POST http://localhost:8081/api/v1/auctions \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Vintage Camera",
    "description": "Rare 1960s camera",
    "start_price": 50,
    "duration_min": 60
  }'
```

### 4. Place a bid

```bash
curl -X POST http://localhost:8081/api/v1/auctions/1/bids \
  -H "Content-Type: application/json" \
  -d '{"bidder_id": "user123", "amount": 75}'
```

### 5. Connect via WebSocket for real-time updates

```javascript
const ws = new WebSocket('ws://localhost:8081/ws/auctions/1');
ws.onmessage = (e) => console.log('New bid:', JSON.parse(e.data));
```

## Configuration

Environment variables (defaults shown):

| Variable | Default | Description |
|----------|---------|-------------|
| PORT | 8080 | Server port |
| DB_HOST | localhost | PostgreSQL host |
| DB_PORT | 5432 | PostgreSQL port |
| DB_USER | auction | Database user |
| DB_PASSWORD | auction | Database password |
| DB_NAME | auction_db | Database name |
| REDIS_ADDR | localhost:6379 | Redis address |
| REDIS_DB | 0 | Redis database number |

## Project Structure

```
.
├── cmd/server/          # Application entry point
├── config/              # Configuration
├── internal/
│   ├── db/              # Database migrations
│   ├── handler/         # HTTP & WebSocket handlers
│   ├── models/          # Data models
│   ├── redis/           # Redis Pub/Sub
│   ├── repository/      # Data access layer
│   └── worker/          # Background workers
├── scripts/             # Test & utility scripts
├── docker-compose.yml
└── go.mod
```

## Testing

### Unit tests

```bash
go test ./...
```

### Integration tests

1. Start services: `docker-compose up -d`
2. Start server: `go run ./cmd/server`
3. Run: `powershell -File ./scripts/integration_test.ps1`

## GitHub Setup

```bash
git init
git add .
git commit -m "Initial commit: Real-Time Auction System"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/real-time-auction.git
git push -u origin main
```

## License

MIT
