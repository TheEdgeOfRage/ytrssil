package db

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"

	ytrssilConfig "gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/config"
)

type postgresDB struct {
	l  *slog.Logger
	db *sql.DB
}

func NewPostgresDB(log *slog.Logger, dbCfg ytrssilConfig.DB) (*postgresDB, error) {
	db, err := sql.Open("postgres", dbCfg.DBURI)
	if err != nil {
		return nil, err
	}

	return &postgresDB{
		l:  log,
		db: db,
	}, nil
}
