.PHONY: migrate-up migrate-down migrate-status migrate-create migrate-reset graphql dev server debug readme seed seed-avatars convert-avatars test test-db-setup setup

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

dev:
	@if command -v air > /dev/null; then \
		air; \
	elif [ -f ~/go/bin/air ]; then \
		~/go/bin/air; \
	else \
		echo "air not found. Install it with: go install github.com/air-verse/air@latest"; \
		exit 1; \
	fi

server: dev

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

seed:
	go run cmd/seed/main.go

setup: migrate-reset migrate-up seed
	@echo "Database setup completed!"

test-db-setup:
	@echo "Setting up test database..."
	@PGPASSWORD=$${DATABASE_PASSWORD:-postgres} psql -h $${DATABASE_HOST:-localhost} -p $${DATABASE_PORT:-5433} -U $${DATABASE_USER:-postgres} -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'moonshine_test'" | grep -q 1 || \
		PGPASSWORD=$${DATABASE_PASSWORD:-postgres} psql -h $${DATABASE_HOST:-localhost} -p $${DATABASE_PORT:-5433} -U $${DATABASE_USER:-postgres} -d postgres -c "CREATE DATABASE moonshine_test"
	@echo "Applying migrations to test database..."
	@DATABASE_NAME=moonshine_test go run cmd/migrate/main.go -command up

test: test-db-setup
	@go test ./... -v 2>&1 | tee /tmp/test_output.txt | awk ' \
	BEGIN { main_pass=0; main_fail=0; sub_pass=0; sub_fail=0 } \
	/^--- PASS/ { main_pass++ } \
	/^--- FAIL/ { main_fail++ } \
	/^    --- PASS/ { sub_pass++ } \
	/^    --- FAIL/ { sub_fail++ } \
	END { \
		total_pass = main_pass + sub_pass; \
		total_fail = main_fail + sub_fail; \
		print ""; \
		print "=== Статистика тестов ==="; \
		print ""; \
		print "Основных тестов:"; \
		print "  ✓ Пройдено: " main_pass; \
		print "  ✗ Провалено: " main_fail; \
		print ""; \
		print "Подтестов:"; \
		print "  ✓ Пройдено: " sub_pass; \
		print "  ✗ Провалено: " sub_fail; \
		print ""; \
		print "Всего:"; \
		print "  ✓ Пройдено: " total_pass; \
		print "  ✗ Провалено: " total_fail; \
		print ""; \
		if (total_fail > 0) exit 1 \
	}'

convert-avatars:
	@if command -v convert > /dev/null || command -v magick > /dev/null; then \
		cd frontend/assets/images/players/avatars && \
		counter=1 && \
		for file in *.gif; do \
			if [ -f "$$file" ]; then \
				if command -v convert > /dev/null; then \
					convert "$$file" "$$counter.png" && \
					echo "Converted $$file to $$counter.png"; \
				elif command -v magick > /dev/null; then \
					magick "$$file" "$$counter.png" && \
					echo "Converted $$file to $$counter.png"; \
				fi && \
				counter=$$((counter + 1)); \
			fi; \
		done && \
		echo "Conversion complete. You can now delete .gif files if needed."; \
	else \
		echo "ImageMagick not found. Install it first:"; \
		echo "  Ubuntu/Debian: sudo apt-get install imagemagick"; \
		echo "  macOS: brew install imagemagick"; \
		echo "  Or convert manually using online tools"; \
		exit 1; \
	fi
