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

## Requirements

- Go 1.24+
- Node.js 18+ and npm
- Docker and Docker Compose

## Services

- **PostgreSQL** (port 5433): Main database
- **ClickHouse** (ports 8123, 9000): Analytics database

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
