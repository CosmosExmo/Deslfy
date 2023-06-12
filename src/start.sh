#!/bin/sh

set -e

echo "RUN DB Migration"
source /app/app.env
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "START THE APP"
exec "$@"