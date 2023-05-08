#!/bin/sh

set -e

echo "running DB migration"
source /app/app.env
echo "$DB_SOURCE"
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"