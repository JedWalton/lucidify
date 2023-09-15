## 

- Run GH actions locally with https://github.com/nektos/act
    - `$ gh act -v`

- Start dev environment.
    - `$ cd docker && docker compose up`
    - if live reloading go doesn't work try:
        - `$ docker compose build --no-cache`
- Migrations
    - `$ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
    - `$ export POSTGRESQL_URL=postgres://postgres:mysecretpassword@localhost:5432/devdb?sslmode=disable`
    - https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md
    - `$ sudo apt-get update & sudo apt-get install postgresql-client`
    - `$ migrate -database ${POSTGRESQL_URL} -path db/migrations up`
    - `$ psql -h localhost -U postgres -d devdb -p 5432 -c "\d confessions"`

