package usecase

import (
	"context"
	"github.com/BigBullas/TP_DB_project/internal/models"
	"github.com/BigBullas/TP_DB_project/internal/pkg/user"
)

type UseCase struct {
	repo user.Repository
}

func NewRepoUseCase(repo user.Repository) user.UseCase { // почему не *user.Repository
	return &UseCase{repo: repo}
}

func (u *UseCase) CreateUser(ctx context.Context, user models.User) ([]models.User, error) {
	usersWithSameInfo, _ := u.repo.CheckUserForUniq(ctx, user)
	if len(usersWithSameInfo) > 0 {
		return usersWithSameInfo, models.Conflict
	}
	err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return []models.User{user}, nil
}

func (u *UseCase) GetUser(ctx context.Context, nickname string) (models.User, error) {
	return u.repo.GetUser(ctx, nickname)
}
