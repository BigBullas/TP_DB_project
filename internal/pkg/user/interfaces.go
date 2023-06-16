package user

import (
	"awesomeProject1/internal/models"
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, user models.User) error
}
