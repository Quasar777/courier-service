# Загруажем переменные из .env файла
include .env
export ${shell sed 's/=.*//' .env}

GOOSE_DRIVER = postgres
GOOSE_DBSTRING = ${DB_CONNECTION_STRING}
MIGRATIONS_DIR = migrations

.PHONY: migrate-up migrate-down migrate-status migrate-create migrate-version run

# Применение миграций
migrate-up:
	@echo "Applying migrations from ${MIGRATIONS_DIR}..."
	GOOSE_DRIVER=${GOOSE_DRIVER} GOOSE_DBSTRING=${GOOSE_DBSTRING} goose -dir ${MIGRATIONS_DIR} up

migrate-up-test:
	@echo "Applying migrations from ${MIGRATIONS_DIR}..."
	GOOSE_DRIVER=${GOOSE_DRIVER} GOOSE_DBSTRING=${GOOSE_DBSTRING} goose -dir ${MIGRATIONS_DIR} up

migrate-down:
	@echo "Rolling back migrations..."
	GOOSE_DRIVER=${GOOSE_DRIVER} GOOSE_DBSTRING=${GOOSE_DBSTRING} goose -dir ${MIGRATIONS_DIR} down

# Просмотр статуса миграций
migrate-status:
	@echo "Migration status in ${MIGRATIONS_DIR}:"
	GOOSE_DRIVER=${GOOSE_DRIVER} GOOSE_DBSTRING=${GOOSE_DBSTRING} goose -dir ${MIGRATIONS_DIR} status