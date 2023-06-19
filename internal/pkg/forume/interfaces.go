package forume

import (
	"context"
	"github.com/BigBullas/TP_DB_project/internal/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user models.User) error
	CheckUserForUniq(ctx context.Context, user models.User) ([]models.User, error)
	GetUser(ctx context.Context, nickname string) (models.User, error)
	ChangeUserInfo(ctx context.Context, user models.User) (models.User, int)
	CreateForum(ctx context.Context, forum models.Forum) ([]models.Forum, int)
	CheckForumForUniq(ctx context.Context, forum models.Forum) ([]models.Forum, int)
	GetForumDetails(ctx context.Context, slug string) (models.Forum, error)
}

type UseCase interface {
	CreateUser(ctx context.Context, user models.User) ([]models.User, error)
	GetUser(ctx context.Context, nickname string) (models.User, error)
	ChangeUserInfo(ctx context.Context, user models.User) (models.User, int)
	CreateForum(ctx context.Context, forum models.Forum) ([]models.Forum, int)
	GetForumDetails(ctx context.Context, slug string) (models.Forum, error)
}
