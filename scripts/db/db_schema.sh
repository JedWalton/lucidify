#!/bin/bash

export PGPASSWORD='mysecretpassword'
export POSTGRESQL_URL=postgres://postgres:mysecretpassword@localhost:5432/devdb?sslmode=disable

# Display all tables
psql -h localhost -U postgres -d devdb -p 5432 -c "\dt"

# Uncomment the following line if you want to display the entire schema
# psql -h localhost -U postgres -d devdb -p 5432 -c "\d"

