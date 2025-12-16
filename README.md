# Moonshine

Backend API server built with Go using Echo framework and PostgreSQL.

## Requirements

- Go 1.24+
- Docker and Docker Compose (for PostgreSQL)
- PostgreSQL 18.1+ (or use Docker Compose)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd moonshine
```

2. Install dependencies:
```bash
go mod download
```

3. Create `.env` file (you can copy from `.env.example` if it exists):
```bash
cp .env.example .env
```

4. Configure environment variables in `.env`:
```env
# Database
DATABASE_HOST=localhost
DATABASE_PORT=5433
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=moonshine
DATABASE_SSL_MODE=disable

# Server
HTTP_ADDR=:8080

# JWT
JWT_KEY=your-secret-key-here
```

## Database Setup

Start PostgreSQL using Docker Compose:
```bash
docker-compose up -d
```

Check that the database is running:
```bash
docker-compose ps
```

## Database Migrations

The project uses [goose](https://github.com/pressly/goose) for database migration management. All migrations are SQL files located in the `migrations/` directory.

**Important:** Migrations must be run manually using goose commands. The application does not run migrations automatically on startup.

### Apply all migrations

```bash
make migrate-up
```

or directly via command:
```bash
go run cmd/migrate/main.go -command up
```

### Rollback last migration

```bash
make migrate-down
```

or:
```bash
go run cmd/migrate/main.go -command down
```

### Show migration status

```bash
make migrate-status
```

or:
```bash
go run cmd/migrate/main.go -command status
```

### Rollback all migrations

```bash
make migrate-reset
```

or:
```bash
go run cmd/migrate/main.go -command down-to 0
```

### Create new migration

```bash
make migrate-create NAME=add_new_field
```

or:
```bash
go run cmd/migrate/main.go -command create add_new_field
```

This will create two files in the `migrations/` folder:
- `YYYYMMDDHHMMSS_add_new_field.up.sql` - up migration
- `YYYYMMDDHHMMSS_add_new_field.down.sql` - down migration

### Available Makefile commands

- `make migrate-up` - apply all migrations
- `make migrate-down` - rollback last migration
- `make migrate-status` - show migration status
- `make migrate-create NAME=name` - create new migration
- `make migrate-reset` - rollback all migrations

## Running the Server

```bash
go run cmd/server/main.go
```

The server will start on the address specified in `HTTP_ADDR` (default `:8080`).

## API Endpoints

### Public routes

- `GET /health` - health check endpoint
- `POST /signup` - register new user
- `POST /signin` - sign in (returns JWT token)

### Protected routes (require JWT token)

- `GET /user` - get current user information

To access protected routes, you need to pass JWT token in the header:
```
Authorization: Bearer <your-jwt-token>
```

## Project Structure

```
moonshine/
├── cmd/
│   ├── server/          # Main server application
│   │   └── main.go
│   └── migrate/         # Migration management command
│       └── main.go
├── internal/
│   ├── api/             # HTTP layer
│   │   ├── handler/      # HTTP handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── auth_handler_test.go
│   │   │   ├── user_handler.go
│   │   │   └── user_handler_test.go
│   │   └── routes.go     # Route configuration
│   ├── domain/           # Domain models (entities)
│   ├── repository/       # Data access layer
│   │   ├── database.go   # Database connection
│   │   ├── database_test.go
│   │   ├── user.go       # User repository
│   │   └── user_test.go
│   ├── service/          # Business logic layer
│   │   ├── user.go
│   │   └── user_test.go
│   └── util/             # Utility functions
│       └── main.go
├── migrations/           # SQL migrations for goose
├── modules/
│   └── database/        # Legacy database connection (old GORM)
├── docker-compose.yml    # PostgreSQL configuration
└── Makefile             # Make commands
```

## Testing

Tests use a separate test database to avoid affecting development data. The test database name is configured via environment variables.

### Test Database Setup

1. Create `.env.test` file (you can copy from `.env.test.example`):
```bash
cp .env.test.example .env.test
```

2. Configure test database in `.env.test`:
```env
DATABASE_HOST=localhost
DATABASE_PORT=5433
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=moonshine_test
DATABASE_SSL_MODE=disable
```

3. Create test database:
```bash
docker-compose exec postgres psql -U postgres -c "CREATE DATABASE moonshine_test;"
```

4. Run migrations on test database:
```bash
DATABASE_NAME=moonshine_test go run cmd/migrate/main.go -command up
```

### Running Tests

Run all tests:
```bash
go test ./...
```

Run tests for a specific package:
```bash
go test ./internal/api/handler
go test ./internal/repository
go test ./internal/service
```

Run tests with verbose output:
```bash
go test -v ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

For UI tests:
```bash
sudo $GOPATH/bin/goconvey
```

**Note:** In Go, tests are located next to the code they test in `*_test.go` files. This is the idiomatic Go way. Tests automatically load `.env.test` if available, otherwise they use environment variables or defaults.

## Development

### Adding new routes

1. Create a handler in `internal/api/handler/`
2. Add the route in `internal/api/routes.go`

### Creating new migrations

1. Create migration: `make migrate-create NAME=description`
2. Fill SQL in created `.up.sql` and `.down.sql` files
3. Apply migration: `make migrate-up`

### Architecture

The project follows a layered architecture:

- **Domain Layer** (`internal/domain/`): Domain models and entities
- **Repository Layer** (`internal/repository/`): Data access and database operations
- **Service Layer** (`internal/service/`): Business logic
- **API Layer** (`internal/api/`): HTTP handlers and routing

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_HOST` | Database host | `localhost` |
| `DATABASE_PORT` | Database port | `5433` |
| `DATABASE_USER` | Database user | `postgres` |
| `DATABASE_PASSWORD` | Database password | `postgres` |
| `DATABASE_NAME` | Database name | `moonshine` |
| `DATABASE_SSL_MODE` | SSL mode | `disable` |
| `HTTP_ADDR` | HTTP server address | `:8080` |
| `JWT_KEY` | JWT secret key | - |

## License

[Specify license]
