#!/bin/bash

go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

sudo apt-get update & sudo apt-get install postgresql-client
