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

## run/test/e2e: run e2e tests alone
.PHONY: run/test/e2e
run/test/e2e:
	@echo "Running e2e tests"
	APP_ENV=test go test -race -vet=off ./cmd/api/tests/...

## run/test/all: run all tests through test files present in the codebase
.PHONY: run/test/all
run/test/all:
	@echo "Running all tests in the codebase"
	@go test -race -vet=off ./...

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

## db/migration/version: get current migration version
.PHONY: db/migration/version
db/migration/version:
	@echo "Retrieving current migration version"
	@migrate -database ${PG_URL} -path ./migrations version

## db/migration/rollback v=$1: rollback to a specific version of the database migration
.PHONY: db/migration/rollback
db/migration/rollback: confirm
	@echo "migrating database to version ${v}..."
	@migrate -database ${PG_URL} -path ./migrations force $(v)

## db/migration/up: apply all up migrations
.PHONY: db/migration/up
db/migration/up: confirm
	@echo "applying all up migrations..."
	@migrate -database ${PG_URL} -path ./migrations up

## db/migration/down: apply all down migrations
.PHONY: db/migration/down
db/migration/down: confirm
	@echo "applying all down migrations..."
	@migrate -database ${PG_URL} -path ./migrations down

## db/migration/crete name=$name: create a new migration file set named $name
.PHONY: db/migration/create
db/migration/create: confirm
	@echo "Creating migration file ${name}"
	@migrate create -ext sql -dir ./migrations -seq $(name)