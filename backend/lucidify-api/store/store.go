package store

import (
	"database/sql"
	"log"
	"os"
)

// Datastore outlines methods that any data storage mechanism should implement.
type Datastore interface {
	Exec(query string, args ...interface{}) (Result, error)
	QueryRow(query string, args ...interface{}) Row
	Close() error
}

// Result represents the result of a database operation.
type Result interface {
	LastInsertId() (int64, error)
}

// Row represents a single row returned from a database query.
type Row interface {
	Scan(dest ...interface{}) error
}

// SQLDB is a concrete implementation of the Datastore interface for a SQL database.
type SQLDB struct {
	*sql.DB
}

func (db *SQLDB) Exec(query string, args ...interface{}) (Result, error) {
	return db.DB.Exec(query, args...)
}

func (db *SQLDB) QueryRow(query string, args ...interface{}) Row {
	return db.DB.QueryRow(query, args...)
}

// DBStore represents our main datastore, containing a reference to an actual or mock database.
type DBStore struct {
	Datastore
	Storer
}

func (store *DBStore) Exec(query string, args ...interface{}) (Result, error) {
	return store.Datastore.Exec(query, args...)
}

func (store *DBStore) QueryRow(query string, args ...interface{}) Row {
	return store.Datastore.QueryRow(query, args...)
}

func (store *DBStore) Close() error {
	return store.Datastore.Close()
}

type Store struct {
	db *sql.DB
}

func NewDBStore(connectionString string) *DBStore {
	rawDB, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	if err := rawDB.Ping(); err != nil {
		log.Fatal(err)
	}

	mainStore := &Store{db: rawDB}
	// storerConfessions := &StorerConfessions{db: rawDB}
	// Initialize other domain-specific stores here if needed...

	return &DBStore{
		Datastore: &SQLDB{DB: rawDB},
		Storer:    mainStore,
	}

}

// ConnectToPostgres initializes the database using a connection string from environment variables.
func ConnectToPostgres() *DBStore {
	connectionString := os.Getenv("POSTGRESQL_URL")
	if connectionString == "" {
		log.Fatal("POSTGRESQL_URL environment variable is not set")
	}
	return NewDBStore(connectionString)
}
