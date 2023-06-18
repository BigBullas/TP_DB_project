package usecase

import (
	"context"
	"github.com/BigBullas/TP_DB_project/internal/models"
	"github.com/BigBullas/TP_DB_project/internal/pkg/forume"
	"net/http"
)

type UseCase struct {
	repo forume.Repository
}

func NewRepoUseCase(repo forume.Repository) forume.UseCase { // почему не *forume.Repository
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

func (u *UseCase) ChangeUserInfo(ctx context.Context, user models.User) (models.User, int) {
	thisUser, err := u.repo.GetUser(ctx, user.NickName)
	if err != nil {
		return models.User{}, http.StatusInternalServerError
	}
	if thisUser == (models.User{}) {
		return models.User{}, http.StatusNotFound
	}
	usersWithSameInfo, _ := u.repo.CheckUserForUniq(ctx, user)
	if len(usersWithSameInfo) > 1 {
		return models.User{}, http.StatusConflict
	}
	return u.repo.ChangeUserInfo(ctx, user)
}