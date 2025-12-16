.PHONY: migrate-up migrate-down migrate-status migrate-create migrate-reset graphql dev debug readme

migrate-up:
	go run cmd/migrate/main.go -command up

migrate-down:
	go run cmd/migrate/main.go -command down

migrate-status:
	go run cmd/migrate/main.go -command status

migrate-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	go run cmd/migrate/main.go -command create $(NAME)

migrate-reset:
	go run cmd/migrate/main.go -command down-to 0

graphql:
	go run github.com/99designs/gqlgen@latest generate --config internal/graphql/gqlgen.yml
	@rm -f internal/graphql/schema.resolvers.go

dev:
	@if command -v air > /dev/null; then \
		air; \
	elif [ -f ~/go/bin/air ]; then \
		~/go/bin/air; \
	else \
		echo "air not found. Install it with: go install github.com/air-verse/air@latest"; \
		exit 1; \
	fi

debug:
	@if command -v dlv > /dev/null; then \
		dlv debug ./cmd/server --headless --listen=:2345 --api-version=2 --accept-multiclient; \
	elif [ -f ~/go/bin/dlv ]; then \
		~/go/bin/dlv debug ./cmd/server --headless --listen=:2345 --api-version=2 --accept-multiclient; \
	else \
		echo "delve not found. Install it with: go install github.com/go-delve/delve/cmd/dlv@latest"; \
		exit 1; \
	fi

readme:
	@if command -v glow > /dev/null; then \
		glow README.md; \
	elif [ -f ~/go/bin/glow ]; then \
		~/go/bin/glow README.md; \
	else \
		echo "glow not found. Install it with: go install github.com/charmbracelet/glow@latest"; \
		echo "Or use VS Code: Press Ctrl+Shift+V to preview markdown"; \
		exit 1; \
	fi
