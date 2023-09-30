package postgresqlclient

import (
	"database/sql"
	"fmt"
	"lucidify-api/modules/config"

	_ "github.com/lib/pq"
)

type PostgreSQL struct {
	db *sql.DB
}

func NewPostgreSQL(postgresqlURL string) (*PostgreSQL, error) {
	config := config.NewServerConfig()
	postgresqlURLFromConfig := config.PostgresqlURL
	if postgresqlURL == "" {
		return nil, fmt.Errorf("POSTGRESQL_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", postgresqlURLFromConfig)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgreSQL{db: db}, nil
}
