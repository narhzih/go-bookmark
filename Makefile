include .env.staging

PG_URL ?= postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_DB_HOST}:${POSTGRES_DB_PORT}/${POSTGRES_DB}?sslmode=disable
MIGRATIONS_PATH ?= $(shell pwd)/sql

migrate-version:
	migrate -database ${PG_URL} -path ./migrations version
migrate-rollback:
	migrate -database ${PG_URL} -path ./migrations force $(v)
migrate-up:
	migrate -database ${PG_URL} -path ./migrations up
migrate-down:
	migrate -database ${PG_URL} -path ./migrations down
# creates a new migration in the sql directory
# e.g: `make migration name=create_users_table` creates a new migration named "create_users_table"
migration:
	migrate create -ext sql -dir ./migrations -seq $(name)