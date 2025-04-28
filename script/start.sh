#!bin/bash

until pg_isready -h db -p 5432 -U task-manager --dbname=task-store; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 1
done

echo "PostgreSQL is ready, starting the app..."

/usr/src/app/task
