#!/bin/bash

export PGPASSWORD='mysecretpassword'
export POSTGRESQL_URL=postgres://postgres:mysecretpassword@localhost:5432/devdb?sslmode=disable

watch "psql -h localhost -U postgres -d devdb -p 5432 -c 'SELECT * FROM users LIMIT 10;'"
