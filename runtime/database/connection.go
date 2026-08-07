package database

import (
	"database/sql"
	"fmt"
)

func Connect(cfg Config) (*Database, error) {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sql.Open(
		cfg.Driver,
		dsn,
	)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return New(db), nil
}