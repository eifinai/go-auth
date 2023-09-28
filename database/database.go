package database

import (
	"database/sql"
)

func NewDatabase(psqlInfo string) (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// func (d Database) Exec(query string) (sql.Result, error) {
// 	return d.DB.Exec(query)
// }

// func (d Database) QueryRow(query string, param string) *sql.Row {
// 	return d.DB.QueryRow(query, param)
// }
