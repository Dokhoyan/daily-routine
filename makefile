include .env
LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.20.0


local-migration-status:
	${LOCAL_BIN}/goose -dir ${MIGRATION_DIR} postgres "host=localhost port=${PG_PORT} dbname=${PG_DATABASE_NAME} user=${PG_USER} password=${PG_PASSWORD} sslmode=disable" status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${MIGRATION_DIR} postgres "host=localhost port=${PG_PORT} dbname=${PG_DATABASE_NAME} user=${PG_USER} password=${PG_PASSWORD} sslmode=disable" up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${MIGRATION_DIR} postgres "host=localhost port=${PG_PORT} dbname=${PG_DATABASE_NAME} user=${PG_USER} password=${PG_PASSWORD} sslmode=disable" down -v


prod-migration-status:
	${LOCAL_BIN}/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} status -v

prod-migration-up:
	${LOCAL_BIN}/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v

prod-migration-down:
	${LOCAL_BIN}/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} down -v


make-migrate:
	bin/goose -dir ./migrations create $(name) sql

run:
	go run cmd/server/main.go --config-path=config/prod/.env  