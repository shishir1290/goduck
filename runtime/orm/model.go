package orm

import "github.com/shishir1290/goduck/runtime/database"

type Model struct {
	db *database.Database
}

func New(db *database.Database) *Model {
	return &Model{
		db: db,
	}
}