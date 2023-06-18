package user

import (
	"context"
	"github.com/BigBullas/TP_DB_project/internal/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user models.User) error
	CheckUserForUniq(ctx context.Context, user models.User) ([]models.User, error)
}

type UseCase interface {
	CreateUser(ctx context.Context, user models.User) ([]models.User, error)
}
