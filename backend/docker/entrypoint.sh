#!/bin/sh
set -eu

run_migrations="${DB_RUN_MIGRATIONS:-true}"
max_attempts="${DB_MIGRATION_MAX_ATTEMPTS:-30}"
retry_delay_seconds="${DB_MIGRATION_RETRY_DELAY_SECONDS:-2}"
migration_allow_dirty="${DB_MIGRATION_ALLOW_DIRTY:-true}"
migration_baseline="${DB_MIGRATION_BASELINE:-}"

if [ -z "${ATLAS_DATABASE_URL:-}" ]; then
  db_host="${DB_HOST:-localhost}"
  db_port="${DB_PORT:-5432}"
  db_user="${DB_USER:-postgres}"
  db_password="${DB_PASSWORD:-postgres}"
  db_name="${DB_DATABASE:-${DB_NAME:-postgres}}"
  db_sslmode="${DB_SSL_MODE:-require}"
  db_timezone="${DB_TIMEZONE:-UTC}"

  export ATLAS_DATABASE_URL="postgresql://${db_user}:${db_password}@${db_host}:${db_port}/${db_name}?sslmode=${db_sslmode}&timezone=${db_timezone}"
fi

if [ "$run_migrations" = "true" ]; then
  echo "Running database migrations..."

  set -- atlas migrate apply --env runtime --url "$ATLAS_DATABASE_URL"

  if [ "$migration_allow_dirty" = "true" ]; then
    set -- "$@" --allow-dirty
  fi

  if [ -n "$migration_baseline" ]; then
    set -- "$@" --baseline "$migration_baseline"
  fi

  attempt=1
  while [ "$attempt" -le "$max_attempts" ]; do
    if "$@"; then
      echo "Migrations applied successfully."
      break
    fi

    echo "Migration attempt ${attempt}/${max_attempts} failed. Retrying in ${retry_delay_seconds}s..."
    attempt=$((attempt + 1))
    sleep "$retry_delay_seconds"
  done

  if [ "$attempt" -gt "$max_attempts" ]; then
    echo "Migration failed after ${max_attempts} attempts. Exiting."
    exit 1
  fi
else
  echo "DB_RUN_MIGRATIONS=false, skipping migrations."
fi

exec /app/server
