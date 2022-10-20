include .env

PG_URL ?= postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_DB_HOST}:${POSTGRES_DB_PORT}/${POSTGRES_DB}?sslmode=${DB_SSL_MODE}
MIGRATIONS_PATH ?= $(shell pwd)/sql

run/api:
	@echo "Running api..."
	APP_ENV=dev go run ./cmd/api

run/build/prod:
	@echo "Running api through local build..."
	APP_ENV=prod ./bin/api

run/build/dev:
	@echo "Running api through local build..."
	APP_ENV=dev ./bin/api

build/api:
	@echo "Creating binary build for ./cmd/api..."
	go build -o=./bin/api ./cmd/api

migrate-version:
	migrate -database ${PG_URL} -path ./migrations version
migrate-rollback:
	migrate -database ${PG_URL} -path ./migrations force $(v)
migrate-up:
	migrate -database ${PG_URL} -path ./migrations up
migrate-down:
	migrate -database ${PG_URL} -path ./migrations down
migrate-drop:
	migrate -database ${PG_URL} -path ./migrations drop -f
# creates a new migration in the sql directory
# e.g: `make migration name=create_users_table` creates a new migration named "create_users_table"
migration:
	migrate create -ext sql -dir ./migrations -seq $(name)