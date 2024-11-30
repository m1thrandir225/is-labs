package db

import "database/sql"

type Store interface {
	Querier
}

type SQLiteStore struct {
	db *sql.DB
	*Queries
}

func NewSQLiteStore(db *sql.DB) Store {
	return &SQLiteStore{
		db:      db,
		Queries: New(db),
	}
}
