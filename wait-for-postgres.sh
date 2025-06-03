#!/bin/sh

until pg_isready -h "$DB_HOST" -p "$DB_PORT" > /dev/null 2>&1; do
  echo "Ждём PostgreSQL..."
  sleep 1
done

echo "PostgreSQL готов, запускаем сервис"
exec ./pastebin
