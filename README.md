# Real-Time Auction Bidding System

A **high-performance, distributed auction platform** built in Go, designed to handle **real-time bids at scale**. The system leverages WebSockets, Redis Pub/Sub, PostgreSQL, and Kubernetes to ensure low-latency bid propagation, high availability, and concurrency correctness.

## Role & Contributions
As the lead developer, I:

- Architected WebSocket auction rooms and asynchronous background workers for real-time bid propagation.  
- Integrated Redis Pub/Sub to broadcast bids across multiple server instances, handling **>10K bids/min** while maintaining <200ms propagation latency.  
- Ensured correctness under concurrency with event ordering, de-duplication, and idempotent database writes.  
- Implemented autoscaling, readiness/liveness probes, and observability dashboards (Prometheus/Grafana) for safe canary rollouts.  
- Followed cloud-native principles and CI/CD practices to deliver production-ready microservices.


## Features

- **REST API** – Manage auctions and bids via JSON endpoints.  
- **WebSocket Updates** – Real-time bid propagation across multiple clients and server instances.  
- **Redis Pub/Sub** – Efficient broadcasting across distributed servers.  
- **Kubernetes Deployment** – Containerized microservices with Horizontal Pod Autoscaling (HPA).  
- **Background Worker** – Automatic auction expiry and notifications.  
- **Concurrency Safety** – Event ordering, deduplication, and idempotent writes ensure data reliability.


## API Endpoints

| Method | Endpoint | Description |
|--------|---------|-------------|
| POST   | `/api/v1/auctions` | Create a new auction |
| GET    | `/api/v1/auctions/:id` | Retrieve auction details |
| POST   | `/api/v1/auctions/:id/bids` | Place a bid |
| GET    | `/api/v1/auctions/:id/bids` | List all bids for an auction |
| WS     | `/ws/auctions/:id` | Real-time bid updates via WebSocket |


## Tech Stack

- **Languages & Frameworks:** Go 1.21+, Gorilla WebSocket, Gorilla Mux  
- **Data Storage:** PostgreSQL (durable storage), Redis Pub/Sub (distributed messaging)  
- **Infrastructure:** Docker, Kubernetes, Helm, CI/CD pipelines  
- **Monitoring & Observability:** Prometheus metrics, Grafana dashboards  
- **Best Practices:** Cloud-native design, autoscaling, idempotent processing, fault-tolerant distributed architecture


## Quick Start

### Prerequisites
- Go 1.21+  
- Docker & Docker Compose (for PostgreSQL and Redis)

### Start Services
```bash
docker-compose up -d

### 2. Run the server

# Windows
.\scripts\run.ps1

# Linux/Mac
./scripts/run.sh

# Or manually with environment variables
DB_PORT=5434 REDIS_ADDR=localhost:6381 PORT=8081 go run ./cmd/server

### 3. Create an auction

curl -X POST http://localhost:8081/api/v1/auctions \
  -H "Content-Type: application/json" \
  -d '{
        "title": "Vintage Camera",
        "description": "Rare 1960s camera",
        "start_price": 50,
        "duration_min": 60
      }'

### 4. Place a bid

curl -X POST http://localhost:8081/api/v1/auctions/1/bids \
  -H "Content-Type: application/json" \
  -d '{"bidder_id": "user123", "amount": 75}'

### 5. Connect via WebSocket for real-time updates

const ws = new WebSocket('ws://localhost:8081/ws/auctions/1');
ws.onmessage = (e) => console.log('New bid:', JSON.parse(e.data));


## Project Structure

.
├── cmd/server/          # Application entry point
├── config/              # Configuration
├── internal/
│   ├── db/              # Database migrations
│   ├── handler/         # HTTP & WebSocket handlers
│   ├── models/          # Data models
│   ├── redis/           # Redis Pub/Sub logic
│   ├── repository/      # Data access layer
│   └── worker/          # Background jobs (auction expiry)
├── scripts/             # Test & utility scripts
├── docker-compose.yml
└── go.mod

## Testing

### Unit tests

go test ./...


### Integration tests

docker-compose up -d
go run ./cmd/server
powershell -File ./scripts/integration_test.ps1

## GitHub Setup

git init
git add .
git commit -m "Initial commit: Real-Time Auction System"
git branch -M main
git remote add origin https://github.com/YOUR_USERNAME/real-time-auction.git
git push -u origin main


## License

MIT
