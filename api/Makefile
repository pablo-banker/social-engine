
MIGRATIONS_PATH:=common/database/migrations

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

tests:
	@echo "Running tests..."
	go test -v ./... -coverprofile=coverage.out -covermode=atomic 

migrate_create:
	@echo "Creating migration..."
	@migrate create -ext sql -dir "${MIGRATIONS_PATH}" -seq $(NAME)

migrate_up:
	@echo "Migrating up..."
	@migrate -database "${DATABASE_URL}" -path "${MIGRATIONS_PATH}" up

migrate_down:
	@echo "Migrating down..."
	@migrate -database "${DATABASE_URL}" -path "${MIGRATIONS_PATH}" down

migrate_force:
	@echo "Migrating force..."
	@migrate -database "${DATABASE_URL}" -path "${MIGRATIONS_PATH}" force $(version)


generate_docs:
	@echo "Generating docs..."
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init --parseDependency --parseInternal

tests_coverage:
	@echo "Generating test coverage report..."
	go tool cover -html=coverage.out -o coverage.html
	@echo "Test coverage report generated: coverage.html"
	go tool cover -html=coverage.out
