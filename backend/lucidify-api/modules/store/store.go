package store

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore() *Store {
	connectionString := os.Getenv("POSTGRESQL_URL")
	if connectionString == "" {
		log.Fatal("POSTGRESQL_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return &Store{db: db}
}

func (store *Store) UploadDocument(dataSchema map[string]string) (int64, error) {
	title := dataSchema["title"]
	content := dataSchema["content"]
	result, err := store.db.Exec("INSERT INTO documents(title, content) VALUES ($1, $2)", title, content)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// func (store *Store) uploaddocument(title, content string) (int64, error) {
// 	result, err := store.db.Exec("INSERT INTO documents(title, content) VALUES ($1, $2)", title, content)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return result.RowsAffected()
// }
//
// func (store *Store) DeleteByID(id int64) error {
// 	_, err := store.db.Exec("DELETE FROM table_name WHERE id = $1", id)
// 	return err
// }
