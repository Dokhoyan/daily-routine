include .env
LOCAL_BIN:=$(CURDIR)/bin


install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.20.0


LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) \
password=$(PG_PASSWORD) sslmode=disable"


local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

# Production migration commands using PG_DSN from .env (for internal connections)
prod-migration-status:
	${LOCAL_BIN}/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} status -v

prod-migration-up:
	${LOCAL_BIN}/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v

prod-migration-down:
	${LOCAL_BIN}/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} down -v

# Production migration commands for external connections (from local machine)
# Uses external IP address for connecting from outside the cluster
PROD_EXTERNAL_DSN="host=89.111.163.253 port=5432 dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

prod-migration-status-external:
	${LOCAL_BIN}/goose -dir ${MIGRATION_DIR} postgres ${PROD_EXTERNAL_DSN} status -v

prod-migration-up-external:
	${LOCAL_BIN}/goose -dir ${MIGRATION_DIR} postgres ${PROD_EXTERNAL_DSN} up -v

prod-migration-down-external:
	${LOCAL_BIN}/goose -dir ${MIGRATION_DIR} postgres ${PROD_EXTERNAL_DSN} down -v

make-migrate:
	bin/goose -dir ./migrations create $(name) sql

run:
	go run cmd/server/main.go --config-path=config/prod/.env  