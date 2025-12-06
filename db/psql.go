package db

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

type postgresDB struct {
	l  *slog.Logger
	db *sql.DB
}

var _ DB = (*postgresDB)(nil)

func NewPostgresDB(log *slog.Logger, dbURI string) (*postgresDB, error) {
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		return nil, err
	}

	return &postgresDB{
		l:  log,
		db: db,
	}, nil
}
