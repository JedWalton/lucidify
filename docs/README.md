- Run GH actions locally with https://github.com/nektos/act
    - `$ gh act -v`

- Start dev environment.
    - `$ cd docker && chmod +x start.sh && ./start.sh`

- Migrations
    - `$ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`
    - `$ export POSTGRESQL_URL=postgres://postgres:mysecretpassword@localhost:5432/devdb?sslmode=disable`
    - https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md
    - `$ sudo apt-get update & sudo apt-get install postgresql-client`
    - `$ migrate -database ${POSTGRESQL_URL} -path db/migrations up`
    - `$ psql -h localhost -U postgres -d devdb -p 5432 -c "\d <table-name>"`


- Clerk auth -> to expose localhost with ngrok use:
    - 2. Expose Your Local Server:
    If your local server is running on localhost:8080, you can expose it using:
    ngrok http 8080
    This will give you a public URL (e.g., https://abcdef1234.ngrok.io) that forwards to your local server.


It is helpful to set these in your ~/.profile for development
export POSTGRESQL_URL=postgres://postgres:mysecretpassword@localhost:5432/devdb?sslmode=disable
export GOPROXY=https://proxy.golang.org


- Testing Sessions
Hi <@1141364780075077803>. In terms of getting the session token with the Go
SDK -- you can't. That portion of Clerk is all frontend driven. For testing you
want to go Clerk Dashboard -> JWT Templates and created a long lived JWT
template. Once down log into your app's frontend, open a console window and do
`await window.Clerk.session.getToken({ template: '<name>' })` to get that long
lived session. Save that and use for testing.
