package handler

import (
	"context"

	"github.com/alexedwards/argon2id"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

func (h *handler) CreateUser(ctx context.Context, user models.User) error {
	hashedPassword, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	return h.db.CreateUser(ctx, user)
}
