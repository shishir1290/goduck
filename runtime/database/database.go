package database

import "database/sql"

type Database struct {
	db *sql.DB
}

func New(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

func (d *Database) DB() *sql.DB {
	return d.db
}