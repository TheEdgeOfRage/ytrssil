package db

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresDB struct {
	l  *slog.Logger
	db *pgxpool.Pool
}

var _ DB = (*postgresDB)(nil)

func NewPostgresDB(log *slog.Logger, dbURI string) (*postgresDB, error) {
	db, err := pgxpool.New(context.Background(), dbURI)
	if err != nil {
		return nil, err
	}

	return &postgresDB{
		l:  log,
		db: db,
	}, nil
}

func (db *postgresDB) Close() {
	db.db.Close()
}
