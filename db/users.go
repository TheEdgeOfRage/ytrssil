package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/alexedwards/argon2id"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

var authenticateUserQuery = `SELECT password FROM users WHERE username = $1`

func (d *postgresDB) AuthenticateUser(ctx context.Context, user models.User) (bool, error) {
	row := d.db.QueryRowContext(ctx, authenticateUserQuery, user.Username)
	var hashedPassword string
	err := row.Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		d.l.Error("Failed to query user for auth", "error", err)
		return false, err
	}

	match, err := argon2id.ComparePasswordAndHash(user.Password, hashedPassword)
	if err != nil {
		d.l.Error("Failed to hash password", "error", err)
		return false, err
	}

	return match, nil
}

var createUserQuery = `INSERT INTO users (username, password) VALUES ($1, $2) ON CONFLICT DO NOTHING`

func (d *postgresDB) CreateUser(ctx context.Context, user models.User) error {
	resp, err := d.db.ExecContext(ctx, createUserQuery, user.Username, user.Password)
	if err != nil {
		d.l.Error("Failed to create user", "error", err)
		return err
	}

	if affected, _ := resp.RowsAffected(); affected == 0 {
		return ErrUserExists
	}

	return nil
}

var deleteUserQuery = `DELETE FROM users WHERE username = $1`

func (d *postgresDB) DeleteUser(ctx context.Context, username string) error {
	_, err := d.db.ExecContext(ctx, deleteUserQuery, username)
	if err != nil {
		d.l.Error("Failed to delete user", "error", err)
		return err
	}

	return nil
}
