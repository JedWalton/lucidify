package postgresqlclient

import (
	"database/sql"
	"fmt"
	"lucidify-api/server/config"

	_ "github.com/lib/pq"
)

type PostgreSQL struct {
	db *sql.DB
}

func NewPostgreSQL() (*PostgreSQL, error) {
	config := config.NewServerConfig()
	postgresqlURL := config.PostgresqlURL
	if postgresqlURL == "" {
		return nil, fmt.Errorf("POSTGRESQL_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", postgresqlURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgreSQL{db: db}, nil
}
