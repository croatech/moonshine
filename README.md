# Moonshine

Backend API server built with Go using Echo framework and PostgreSQL, with React frontend.

## Quick Start

1. Start PostgreSQL:
```bash
docker-compose up -d
```

2. Setup database (migrations + seed):
```bash
make setup
```

Or step by step:
```bash
make migrate-up
make seed
```

3. Start backend server:
```bash
go run cmd/server/main.go
```

4. Start frontend (in a new terminal):
```bash
cd frontend
npm install
npm run dev
```

The frontend will be available at `http://localhost:3000` and the backend at `http://localhost:8080`.

## Requirements

- Go 1.24+
- Node.js 18+ and npm (for frontend)
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

3. Create `.env` file:
```bash
cp .env.example .env
```

4. Configure environment variables in `.env`:
```env
DATABASE_HOST=localhost
DATABASE_PORT=5433
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=moonshine
DATABASE_SSL_MODE=disable
HTTP_ADDR=:8080
JWT_KEY=your-secret-key-here
ENV=development
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

Migrations must be run manually using goose commands. The application does not run migrations automatically on startup.

### Apply all migrations

```bash
make migrate-up
```

### Rollback last migration

```bash
make migrate-down
```

### Show migration status

```bash
make migrate-status
```

### Rollback all migrations

```bash
make migrate-reset
```

### Create new migration

```bash
make migrate-create NAME=add_new_field
```

This will create two files in the `migrations/` folder:
- `YYYYMMDDHHMMSS_add_new_field.up.sql` - up migration
- `YYYYMMDDHHMMSS_add_new_field.down.sql` - down migration

### Setup database (reset + migrate + seed)

```bash
make setup
```

This command will:
1. Rollback all migrations
2. Apply all migrations
3. Seed the database with initial data

### Available Makefile commands

- `make migrate-up` - apply all migrations
- `make migrate-down` - rollback last migration
- `make migrate-status` - show migration status
- `make migrate-create NAME=name` - create new migration
- `make migrate-reset` - rollback all migrations
- `make setup` - reset, migrate and seed database
- `make seed` - seed database with initial data
- `make graphql` - generate GraphQL code from schema
- `make dev` - run server with hot reload (air)
- `make debug` - run server with Delve debugger
- `make test` - run all tests with statistics
- `make test-db-setup` - setup test database

## Running the Server

### Backend Server

#### Standard Run

```bash
go run cmd/server/main.go
```

The server will start on the address specified in `HTTP_ADDR` (default `:8080`).

#### Hot Reload (Development)

For automatic server restart on file changes, use `air`:

```bash
make dev
```

If `air` is not installed, install it first:
```bash
go install github.com/air-verse/air@latest
```

Configuration is in `.air.toml` file.

### Frontend Server

The frontend is built with React and Vite. To run it:

1. Install dependencies:
```bash
cd frontend
npm install
```

2. Start the development server:
```bash
npm run dev
```

The frontend will start on `http://localhost:3000` and will proxy requests to the backend at `http://localhost:8080`.

Available pages:
- `/signup` - User registration
- `/signin` - User login

Make sure the backend server is running before starting the frontend.

### Debugging with Delve

The project is configured for debugging with [Delve](https://github.com/go-delve/delve), the Go debugger.

#### Installation

Install Delve:
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

#### Method 1: Breakpoints in VS Code (Recommended)

1. Open the file
2. Click left of the line number - a red dot will appear (breakpoint)
3. Press F5 to start debugging
4. Execution will automatically stop at that line

Available debug configurations:
- Debug Server - Launch server in debug mode
- Attach to Server - Attach to running process
- Debug Current Test - Debug current test file

#### Method 2: Breakpoints from Command Line

Start server with Delve:
```bash
make debug
```

#### Useful Delve Commands

- `break <file>:<line>` - Set breakpoint at line
- `break <function>` - Set breakpoint at function
- `breakpoints` - Show all breakpoints
- `clear <id>` - Remove breakpoint
- `clearall` - Remove all breakpoints
- `condition <id> <expr>` - Set condition for breakpoint
- `continue` or `c` - Continue execution
- `next` or `n` - Next line
- `step` or `s` - Step into function
- `stepout` - Step out of function
- `print <var>` or `p <var>` - Print variable value
- `locals` - Show local variables
- `args` - Show function arguments
- `stack` - Show call stack
- `goroutines` - Show all goroutines
- `exit` or `quit` - Exit debugger

## GraphQL API

The project uses GraphQL for all API operations. After starting the server, you can access:

- GraphQL Endpoint: http://localhost:8080/graphql
- GraphQL Schema: http://localhost:8080/schema.graphql

### Using Altair GraphQL Client

Development mode (ENV != production):

1. Open Altair GraphQL Client
2. Set the endpoint URL to: `http://localhost:8080/graphql`
3. Click on "Schema" tab or use "Load Schema" button
4. The schema will be automatically loaded via introspection (no token required)

Alternatively, you can load the schema file directly:
- Schema file URL: `http://localhost:8080/schema.graphql`

Production mode (ENV=production):

- Introspection queries are disabled for security
- Schema endpoint (`/schema.graphql`) is disabled

This prevents attackers from discovering your API structure. Only public operations (`signUp`, `signIn`) are available without authentication.

### Available Queries

- `currentUser` - Get current authenticated user information

### Available Mutations

- `signUp(input: SignUpInput!)` - Register a new user
- `signIn(input: SignInInput!)` - Sign in and get JWT token

### Authentication

For protected queries (like `currentUser`), you need to pass JWT token in the header:
```
Authorization: Bearer <your-jwt-token>
```

### Example Queries

Sign Up:
```graphql
mutation {
  signUp(input: {
    username: "john_doe"
    email: "john@example.com"
    password: "password123"
  }) {
    token
    user {
      id
      username
      email
    }
  }
}
```

Sign In:
```graphql
mutation {
  signIn(input: {
    username: "john_doe"
    password: "password123"
  }) {
    token
    user {
      id
      username
      email
      hp
      level
    }
  }
}
```

Get Current User:
```graphql
query {
  currentUser {
    id
    username
    email
    hp
    level
    gold
    exp
  }
}
```

### Generating GraphQL Code

After modifying the GraphQL schema (`internal/graphql/schema.graphqls`), regenerate the code:

```bash
make graphql
```

This will update the generated files in `internal/graphql/generated/` and `internal/graphql/models/`.

## API Endpoints

### Public routes

- `GET /health` - health check endpoint
- `POST /graphql` - GraphQL endpoint (public mutations like signUp, signIn)

### Protected routes (require JWT token)

- `POST /graphql` - GraphQL endpoint (protected queries like currentUser)

## Project Structure

```
moonshine/
├── cmd/
│   ├── server/          # Main server application
│   │   └── main.go
│   ├── migrate/         # Migration management command
│   │   └── main.go
│   └── seed/            # Database seeding command
│       └── main.go
├── internal/
│   ├── api/             # HTTP layer
│   │   ├── handlers/    # HTTP handlers
│   │   ├── services/    # Business logic services
│   │   ├── middleware/  # HTTP middleware
│   │   └── routes.go    # Route configuration
│   ├── domain/          # Domain models (entities)
│   ├── repository/      # Data access layer
│   ├── worker/          # Background workers
│   └── util/            # Utility functions
├── migrations/          # SQL migrations for goose
├── frontend/            # React frontend
├── docker-compose.yml   # PostgreSQL configuration
└── Makefile            # Make commands
```

## Testing

Tests use a separate test database (`moonshine_test`) to avoid affecting development data.

### Test Database Setup

The test database is automatically created and migrated when running `make test`. For manual setup:

1. Create `.env.test` file:
```bash
cat > .env.test << EOF
DATABASE_HOST=localhost
DATABASE_PORT=5433
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=moonshine_test
DATABASE_SSL_MODE=disable
EOF
```

2. Setup test database manually:
```bash
make test-db-setup
```

This will:
- Create `moonshine_test` database if it doesn't exist
- Apply all migrations to test database

### Running Tests

Run all tests with statistics:
```bash
make test
```

This will:
- Setup test database automatically
- Run all tests
- Display statistics (passed/failed tests)

Run tests for a specific package:
```bash
go test ./internal/repository
go test ./internal/api/services
```

Run tests with verbose output:
```bash
go test -v ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Tests automatically load `.env.test` if available, otherwise they use environment variables or defaults.

## Development

### Development Workflow

Recommended workflow:
1. Use `make dev` for hot reload during development
2. Use Delve debugger (F5 in VS Code) when you need to debug specific issues
3. Use standard `go run` for quick testing

### Adding new GraphQL operations

1. Update the GraphQL schema in `internal/graphql/schema.graphqls`
2. Run `make graphql` to regenerate code
3. Implement the resolver methods in `internal/graphql/query.go` or `internal/graphql/mutation.go`

### Creating new migrations

1. Create migration: `make migrate-create NAME=description`
2. Fill SQL in created `.up.sql` and `.down.sql` files
3. Apply migration: `make migrate-up`

### Architecture

The project follows a layered architecture:

- Domain Layer (`internal/domain/`): Domain models and entities
- Repository Layer (`internal/repository/`): Data access and database operations
- Service Layer (`internal/api/services/`): Business logic
- Handler Layer (`internal/api/handlers/`): HTTP handlers and request/response handling
- API Layer (`internal/api/`): HTTP routing and middleware

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
| `ENV` | Environment (development/production) | `development` |

## License

[Specify license]
