package storage

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	Shutdown() error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	connStr := "user=" + user + " password=" + password + " dbname=" + dbname + " host=" + host + " port=" + port + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(20)

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Shutdown() error {
	return s.db.Close()
}
