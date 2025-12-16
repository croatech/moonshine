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

# Environment (development/production)
# In production, introspection is disabled for security
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
- `make graphql` - generate GraphQL code from schema
- `make dev` - run server with hot reload (air)
- `make debug` - run server with Delve debugger

## Running the Server

### Standard Run

```bash
go run cmd/server/main.go
```

The server will start on the address specified in `HTTP_ADDR` (default `:8080`).

### Hot Reload (Development)

For automatic server restart on file changes, use `air`:

```bash
make dev
```

Or directly:
```bash
air
```

If `air` is not installed, install it first:
```bash
go install github.com/air-verse/air@latest
```

**What it does:**
- Automatically rebuilds the server when `.go` files change
- Restarts the server automatically
- Shows build errors in real-time
- Excludes test files, migrations, and generated code from watching

**Configuration:** Settings are in `.air.toml` file.

### Debugging with Delve

The project is configured for debugging with [Delve](https://github.com/go-delve/delve), the Go debugger.

#### Installation

Install Delve:
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

#### Method 1: Breakpoints in VS Code (Recommended)

1. **Open the file** (e.g., `internal/graphql/query.go`)
2. **Click left of the line number** - a red dot will appear (breakpoint)
3. **Press F5** to start debugging
4. Execution will **automatically stop** at that line

**Example:**
- Open `internal/graphql/query.go`
- Set breakpoint on line 13 (click left of `userID, err := getUserIDFromContext(ctx)`)
- Press F5
- Execute GraphQL query `currentUser`
- Execution will stop at line 13

**Available debug configurations:**
- **Debug Server** - Launch server in debug mode
- **Attach to Server** - Attach to running process
- **Debug Current Test** - Debug current test file

#### Method 2: Breakpoints from Command Line

**Start server with Delve:**
```bash
dlv debug ./cmd/server
```

Or use Makefile:
```bash
make debug
```

**In Delve interactive console:**
```bash
(dlv) break internal/graphql/query.go:13
Breakpoint 1 set at 0x1234567 for moonshine/internal/graphql.(*queryResolver).CurrentUser() ./internal/graphql/query.go:13

(dlv) continue
```

#### Method 3: Programmatic Breakpoint in Code

Add `runtime.Breakpoint()` at the desired location:

```go
import "runtime"

func (r *queryResolver) CurrentUser(ctx context.Context) (*models.User, error) {
    userID, err := getUserIDFromContext(ctx)
    if err != nil {
        return nil, errors.New("unauthorized: invalid or missing token")
    }

    runtime.Breakpoint()  // Will stop here when running under debugger

    user, err := r.userRepo.FindByID(userID)
    // ...
}
```

**Important:** `runtime.Breakpoint()` only works when the program is run under a debugger. If run without a debugger, the program will crash.

#### Method 4: Conditional Breakpoints in VS Code

1. Set a regular breakpoint
2. **Right-click** on the breakpoint → "Edit Breakpoint"
3. Add a condition, e.g.: `userID == 1`

Execution will stop only if the condition is true.

#### Method 5: Conditional Breakpoints in Delve

```bash
(dlv) break internal/graphql/query.go:13
Breakpoint 1 set at 0x1234567

(dlv) condition 1 userID == 1
```

Will stop only if `userID == 1`.

#### Practical Example: Debugging `CurrentUser`

**In VS Code:**
1. Open `internal/graphql/query.go`
2. Set breakpoint on line 13
3. Press F5 (select "Debug Server")
4. In Altair, execute the query:
   ```graphql
   query {
     currentUser {
       id
       username
     }
   }
   ```
5. Execution will stop at line 13
6. In the debugger you can:
   - View `ctx` in Variables panel
   - Press `F10` (next) to move to the next line
   - Press `F11` (step into) to enter `getUserIDFromContext`
   - Hover over variables to see their values

**From command line:**
```bash
$ dlv debug ./cmd/server
Type 'help' for list of commands.
(dlv) break internal/graphql/query.go:13
Breakpoint 1 set at 0x1234567 for moonshine/internal/graphql.(*queryResolver).CurrentUser() ./internal/graphql/query.go:13
(dlv) continue
> moonshine/internal/graphql.(*queryResolver).CurrentUser() ./internal/graphql/query.go:13 (hits goroutine(1):1 total:1) (PC: 0x1234567)
    12:	func (r *queryResolver) CurrentUser(ctx context.Context) (*models.User, error) {
=>  13:		userID, err := getUserIDFromContext(ctx)
    14:		if err != nil {
    15:			return nil, errors.New("unauthorized: invalid or missing token")
    16:		}
(dlv) print ctx
context.Context = ...
(dlv) next
> moonshine/internal/graphql.(*queryResolver).CurrentUser() ./internal/graphql/query.go:18 (PC: 0x1234568)
    18:		user, err := r.userRepo.FindByID(userID)
(dlv) print userID
uint = 1
(dlv) continue
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

#### Tips

1. **Use VS Code** for visual debugging - it's the most convenient way
2. **Conditional breakpoints** are useful for debugging specific cases
3. **Don't leave `runtime.Breakpoint()` in code** - remove it after debugging
4. **Use `print` and `locals`** for quick variable inspection
5. **`step` vs `next`**: `step` enters functions, `next` skips them

## GraphQL API

The project uses GraphQL for all API operations. After starting the server, you can access:

- **GraphQL Endpoint**: http://localhost:8080/graphql
- **GraphQL Schema**: http://localhost:8080/schema.graphql

### Using Altair GraphQL Client

**Development mode (ENV != production):**

1. Open Altair GraphQL Client
2. Set the endpoint URL to: `http://localhost:8080/graphql`
3. Click on "Schema" tab or use "Load Schema" button
4. The schema will be automatically loaded via introspection (no token required)

Alternatively, you can load the schema file directly:
- Schema file URL: `http://localhost:8080/schema.graphql`

**Production mode (ENV=production):**

- Introspection queries are **disabled** for security
- Schema endpoint (`/schema.graphql`) is **disabled**

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

**Sign Up:**
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

**Sign In:**
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

**Get Current User:**
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
│   └── migrate/         # Migration management command
│       └── main.go
├── internal/
│   ├── api/             # HTTP layer
│   │   └── routes.go     # Route configuration
│   ├── graphql/          # GraphQL layer
│   │   ├── schema.graphqls  # GraphQL schema
│   │   ├── resolver.go      # Base resolver
│   │   ├── query.go         # Query resolvers
│   │   ├── mutation.go      # Mutation resolvers
│   │   ├── helpers.go       # Helper functions
│   │   ├── handler.go       # GraphQL handler for Echo
│   │   ├── generated/       # Generated GraphQL code
│   │   └── models/          # Generated models
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
│       └── hash.go
├── migrations/           # SQL migrations for goose
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

### Development Workflow

**Recommended workflow:**
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

- **Domain Layer** (`internal/domain/`): Domain models and entities
- **Repository Layer** (`internal/repository/`): Data access and database operations
- **Service Layer** (`internal/service/`): Business logic
- **GraphQL Layer** (`internal/graphql/`): GraphQL schema, resolvers, and handlers
- **API Layer** (`internal/api/`): HTTP routing and middleware

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
