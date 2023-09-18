// modules/testutils/db.go

package testutils

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func SetupDB() *sql.DB {
	godotenv.Load("../../../../.env")
	postgresqlURL := os.Getenv("POSTGRESQL_URL")

	db, err := sql.Open("postgres", postgresqlURL)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	return db
}
