include .env

# --- HELPERS ---
## help: display this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## help/api: display API usage
.PHONY: help/api
help/api:
	@go run ./cmd/api --help

# confirm: display confirmation prompt
.PHONY: confirm
confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

# --- DEVELOPMENT ---
## build/api/dev: build the cmd/api application in development mode
.PHONY: build/api/dev
build/api/dev:
	@echo "Creating binary build for ./cmd/api..."
	go build -o=./bin/api ./cmd/api

## run/build/dev: run the build in bin/api in development mode
.PHONY: run/build/dev
run/build/dev:
	@echo "Running api through dev build..."
	APP_ENV=dev ./bin/api

## run/api: run the cmd/api application in development mode
.PHONY: run/api
run/api:
	@echo "Running api directly..."
	APP_ENV=dev go run ./cmd/api

# --- TEST COMMANDS ---

## test/api/e2e: run e2e tests alone
.PHONY: test/api/e2e
test/api:
	@echo "Running e2e tests on ./cmd/api..."
	APP_ENV=test go test -race -vet=off ./cmd/api/tests/...

## test/db: run unit tests for database operations
.PHONY: test/db
test/db:
	@echo "Running database tests alone"
	APP_ENV=test go test -race -vet=off ./db/...

# --- QUALITY CONTROL ---
## audit: tidy and vendor dependencies and format, vet and test codebase
.PHONY: audit
audit:
	@echo "formatting codebase..."
	@go fmt ./...

	@echo "vetting code..."
	@go vet ./...
	@staticcheck ./...

	@echo "running tests..."
	@go test -race -vet=off ./...

# --- DATABASE MIGRATIONS ---
PG_URL ?= postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_DB_HOST}:${POSTGRES_DB_PORT}/${POSTGRES_DB}?sslmode=${DB_SSL_MODE}
MIGRATIONS_PATH ?= $(shell pwd)/sql

## db/migrations/version: get current migration version
.PHONY: db/migrations/version
db/migrations/version:
	@echo "Retrieving current migration version"
	@migrate -database ${PG_URL} -path ./migrations version

## db/migrations/rollback v=$1: rollback to a specific version of the database migration
.PHONY: db/migrations/rollback
db/migrations/rollback: confirm
	@echo "migrating database to version ${v}..."
	@migrate -database ${PG_URL} -path ./migrations force $(v)

## db/migrations/up: apply all up migrations
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo "applying all up migrations..."
	@migrate -database ${PG_URL} -path ./migrations up

## db/migrations/down: apply all down migrations
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo "applying all down migrations..."
	@migrate -database ${PG_URL} -path ./migrations down

## db/migrations/crete name=$name: create a new migration file set named $name
.PHONY: db/migrations/create
db/migrations/create: confirm
	@echo "Creating migration file ${name}"
	@migrate create -ext sql -dir ./migrations -seq $(name)