#!/bin/bash

PG_DATABASE_NAME=${PG_DATABASE_NAME:-daily_routine}
PG_USER=${PG_USER:-postgres}
PG_PASSWORD=${PG_PASSWORD:-postgres}
MIGRATION_DIR=${MIGRATION_DIR:-migrations}

export MIGRATION_DSN="host=pg port=5432 dbname=${PG_DATABASE_NAME} user=${PG_USER} password=${PG_PASSWORD} sslmode=disable"

echo "Waiting for PostgreSQL to be ready..."
for i in {1..30}; do
  if PGPASSWORD="${PG_PASSWORD}" psql -h pg -U "${PG_USER}" -d "${PG_DATABASE_NAME}" -c '\q' 2>/dev/null; then
    echo "PostgreSQL is ready!"
    break
  fi
  echo "PostgreSQL is unavailable - sleeping (attempt $i/30)"
  sleep 2
done

if ! PGPASSWORD="${PG_PASSWORD}" psql -h pg -U "${PG_USER}" -d "${PG_DATABASE_NAME}" -c '\q' 2>/dev/null; then
  echo "ERROR: PostgreSQL is not ready after 60 seconds"
  exit 1
fi

echo "PostgreSQL is ready!"
echo "Running migrations from ${MIGRATION_DIR}..."
goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v