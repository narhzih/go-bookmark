include .env

PG_URL ?= postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost${PG_HOST}:${POSTGRES_DB_PORT}/${POSTGRES_DB}?sslmode=disable
MIGRATIONS_PATH ?= $(shell pwd)/sql

migrate-up:
	migrate -database ${PG_URL} -path ./migrations up
migrate-down:
	migrate -database ${PG_URL} -path ./migrations down
# creates a new migration in the sql directory
# e.g: `make migration name=create_users_table` creates a new migration named "create_users_table"
migration:
	migrate create -ext sql -dir ./migrations -seq $(name)