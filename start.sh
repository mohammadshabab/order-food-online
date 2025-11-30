#!/bin/sh
set -e

echo "Waiting for MariaDB..."

# Wait for MariaDB to be ready
until mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" \
    -e "SELECT 1;" >/dev/null 2>&1; do
  echo "MariaDB is unavailable - retrying..."
  sleep 2
done

echo "MariaDB is up! Running migrations..."

# Run migrations
for f in /migrations/*.up.sql; do
  [ -f "$f" ] || continue
  echo "Applying migration: $f"
  mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" < "$f"
done

echo "Migrations complete! Starting API..."
exec ./order-api
