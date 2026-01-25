# Moonshine

Backend API server built with Go using Echo framework, PostgreSQL, and ClickHouse for analytics. React frontend.

## Quick Start

1. Start services:
```bash
docker-compose up -d
```

2. Setup database:
```bash
make setup
```

3. Start backend:
```bash
go run cmd/server/main.go
```

4. Start frontend:
```bash
cd frontend
npm install
npm run dev
```

Frontend: `http://localhost:3000`  
Backend: `http://localhost:8080`

## LAN access (e.g. from phone)

1. Set `HTTP_ADDR=0.0.0.0:8080` in `.env` (or export it). Default in `.env.example`.
2. Start backend and frontend as above. Vite listens on all interfaces (`host: true`).
3. Get your Mac’s IP: `ipconfig getifaddr en0` (often `192.168.x.x`).
4. On the phone (same Wi‑Fi) open `http://192.168.x.x:3000`.

API and WebSocket use relative URLs, so they go through the frontend origin; no extra config.

## Requirements

- Go 1.24+
- Node.js 18+ and npm
- Docker and Docker Compose

## Services

| Service | Port | Purpose |
|---------|------|---------|
| PostgreSQL | 5433 | Main database |
| ClickHouse | 8123, 9000 | Analytics database |
| Prometheus | 9090 | Metrics collection |
| Grafana | 3001 | Monitoring dashboards |
| Loki | 3100 | Log aggregation |
| cAdvisor | 8088 | Container metrics |

## Monitoring

Access:
- Grafana: http://localhost:3001 (admin/admin)
- Prometheus: http://localhost:9090
- API Metrics: http://localhost:1666/metrics
- cAdvisor: http://localhost:8088

### Metrics

- `http_requests_total` - requests/sec by endpoint, method, status
- `http_request_duration_seconds` - latency histogram (p50, p95, p99)
- `active_websocket_connections` - active WS connections
- `moonshine_fights_total` - total fights
- `moonshine_fight_duration_seconds` - fight duration histogram
- `moonshine_players_online` - online players
- `container_cpu_usage_seconds_total` - CPU usage
- `container_memory_usage_bytes` - RAM usage

### Useful PromQL Queries

Top 5 slowest endpoints (p95):
```promql
topk(5, histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])))
```

5xx error rate:
```promql
rate(http_requests_total{status=~"5.."}[5m])
```

Requests per second by endpoint:
```promql
sum(rate(http_requests_total[5m])) by (path)
```

### Logs

In Grafana, go to Explore → Loki:

View container logs:
```logql
{container="moonshine-postgres-1"}
```

Filter errors:
```logql
{container="moonshine-postgres-1"} |= "ERROR"
```

### Adding Custom Metrics

```go
import "moonshine/internal/metrics"

metrics.FightsTotal.Inc()

start := time.Now()
// ... operation ...
metrics.FightDuration.Observe(time.Since(start).Seconds())

metrics.PlayersOnline.Set(float64(count))
```

## ClickHouse Analytics

ClickHouse automatically replicates `movement_logs` and `rounds` tables from PostgreSQL using WAL replication.

### Check replication

```bash
docker-compose exec clickhouse clickhouse-client
```

```sql
SHOW DATABASES;
SHOW TABLES FROM moonshine_analytics;
SELECT * FROM moonshine_analytics.movement_logs LIMIT 10;
SELECT * FROM moonshine_analytics.rounds LIMIT 10;
```

### Query examples

```sql
-- Top visited cells (last hour)
SELECT to_cell, count() as visits
FROM moonshine_analytics.movement_logs
WHERE created_at >= now() - INTERVAL 1 HOUR
GROUP BY to_cell
ORDER BY visits DESC
LIMIT 10;

-- Popular routes
SELECT from_cell, to_cell, count() as count
FROM moonshine_analytics.movement_logs
WHERE from_cell != ''
GROUP BY from_cell, to_cell
ORDER BY count DESC
LIMIT 20;

-- Round statistics
SELECT 
    date(created_at) as day,
    count() as total_rounds,
    countIf(winner_id IS NOT NULL) as finished_rounds
FROM moonshine_analytics.rounds
GROUP BY day
ORDER BY day DESC;
```

### Add tables to replication

Edit `clickhouse-init.sql`:
```sql
materialized_postgresql_tables_list = 'movement_logs,rounds,new_table'
```

Then recreate:
```bash
docker-compose exec clickhouse clickhouse-client -q "DROP DATABASE moonshine_analytics;"
docker-compose restart clickhouse
```

## Database Migrations

Apply all migrations:
```bash
make migrate-up
```

Rollback last:
```bash
make migrate-down
```

Show status:
```bash
make migrate-status
```

Create new:
```bash
make migrate-create NAME=add_new_field
```

Reset and seed:
```bash
make setup
```

## Makefile Commands

- `make migrate-up` - apply migrations
- `make migrate-down` - rollback last migration
- `make migrate-status` - show status
- `make migrate-create NAME=name` - create migration
- `make migrate-reset` - rollback all
- `make setup` - reset + migrate + seed
- `make seed` - seed database
- `make dev` - run with hot reload (air)
- `make debug` - run with Delve debugger
- `make test` - run tests
- `make swagger` - generate Swagger docs

## Hot Reload

```bash
make dev
```

Requires air:
```bash
go install github.com/air-verse/air@latest
```

## Debugging

Install Delve:
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

Run:
```bash
make debug
```

Or use F5 in VS Code.

## API Endpoints

- `GET /health` - health check
- `POST /api/auth/signup` - register
- `POST /api/auth/signin` - login
- `GET /api/users/me` - current user (requires auth)

## Project Structure

```
moonshine/
├── cmd/
│   ├── server/          # Main server
│   ├── migrate/         # Migrations
│   └── seed/            # Seed data
├── internal/
│   ├── api/             # HTTP layer
│   │   ├── handlers/    # Request handlers
│   │   ├── services/    # Business logic
│   │   ├── middleware/  # Auth, CORS, etc
│   │   └── routes.go    # Routes
│   ├── domain/          # Domain models
│   ├── repository/      # Database access
│   ├── worker/          # Background workers
│   └── util/            # Utilities
├── migrations/          # SQL migrations
├── frontend/            # React app
├── docker-compose.yml   # Services config
└── Makefile
```

## Testing

Run all tests:
```bash
make test
```

Test database setup:
```bash
make test-db-setup
```

Tests use separate `moonshine_test` database.

## Environment Variables

Create `.env` from `.env.example`:

```env
APP_PORT=1666
JWT_KEY=secret

DB_HOST=postgres
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=moonshine
DB_SSLMODE=disable
```
