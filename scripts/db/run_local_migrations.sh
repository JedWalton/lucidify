#!/bin/bash

export POSTGRESQL_URL=postgres://postgres:mysecretpassword@localhost:5432/devdb?sslmode=disable

migrate -database ${POSTGRESQL_URL} -path ../../backend/db/migrations up
