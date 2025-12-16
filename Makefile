.PHONY: migrate-up migrate-down migrate-status migrate-create migrate-reset

# Применить все миграции
migrate-up:
	go run cmd/migrate/main.go -command up

# Откатить последнюю миграцию
migrate-down:
	go run cmd/migrate/main.go -command down

# Показать статус миграций
migrate-status:
	go run cmd/migrate/main.go -command status

# Создать новую миграцию (использование: make migrate-create NAME=migration_name)
migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	go run cmd/migrate/main.go -command create $(NAME)

# Откатить все миграции
migrate-reset:
	go run cmd/migrate/main.go -command down-to 0

