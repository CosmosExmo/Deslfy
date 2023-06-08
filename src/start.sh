#!/bin/sh

include src/app.env
export

set -e

echo "RUN DB Migration"

/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "START THE APP"
exec "$@"