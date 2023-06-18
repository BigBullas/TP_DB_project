package usecase

import (
	"awesomeProject1/internal/models"
	"awesomeProject1/internal/pkg/user"
	"context"
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
