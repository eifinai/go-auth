package database

import (
	"database/sql"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase(psqlInfo string) (Database, error) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return Database{}, err
	}

	err = db.Ping()
	if err != nil {
		return Database{}, err
	}

	return Database{DB: db}, nil
}

func (d Database) Close() error {
	return d.DB.Close()
}

func (d Database) Exec(query string) (sql.Result, error) {
	return d.DB.Exec(query)
}

func (d Database) QueryRow(query string) *sql.Row {
	return d.DB.QueryRow(query)
}
