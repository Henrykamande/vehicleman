package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {

	// PostgreSQL connection string
	//postgresql://user:1234@localhost:5432/property?sslmode=disable
	//connStr := "host=localhost port=5432 dbname=property user=postgres password=1234 sslmode=disable "

	db, err := sql.Open("postgres", "postgresql://postgres:1234@localhost:5432/lorry_management?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}
