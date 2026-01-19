# Database Migrations

The project uses [goose](https://github.com/pressly/goose) for database migration management.

## Installation

Goose is already added to project dependencies. If you need to install it separately:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Usage

### Running migrations

Apply all migrations:
```bash
go run cmd/migrate/main.go -command up
```

Rollback last migration:
```bash
go run cmd/migrate/main.go -command down
```

Rollback all migrations:
```bash
go run cmd/migrate/main.go -command down-to 0
```

Check migration status:
```bash
go run cmd/migrate/main.go -command status
```

### Creating new migration

```bash
go run cmd/migrate/main.go -command create migration_name
```

This will create two files:
- `YYYYMMDDHHMMSS_migration_name.up.sql` - up migration
- `YYYYMMDDHHMMSS_migration_name.down.sql` - down migration

## Migration Structure

All migrations are located in the `migrations/` folder and follow the format:
- `00001_create_avatars.sql` - create avatars table
- `00002_create_equipment_categories.sql` - create equipment_categories table
- etc.

Each migration contains:
- `-- +goose Up` - section for applying migration
- `-- +goose Down` - section for rolling back migration

## Environment Variables

Migrations use the same environment variables as the main application:

- `DATABASE_HOST` (default: localhost)
- `DATABASE_PORT` (default: 5433)
- `DATABASE_USER` (default: postgres)
- `DATABASE_PASSWORD` (default: postgres)
- `DATABASE_NAME` (default: moonshine)
- `DATABASE_SSL_MODE` (default: disable)

## Migration Order

Migrations are created in the correct order considering dependencies:

1. avatars
2. equipment_categories
3. tool_categories
4. equipment_items
5. tool_items
6. locations
7. bots
8. users (depends on avatars, locations)
9. location_bots (depends on locations, bots)
10. location_locations (depends on locations)
11. fights (depends on users, bots)
12. rounds (depends on fights)
13. inventory (depends on users, equipment_items)
