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
	CreateThread(ctx context.Context, thread models.Thread) ([]models.Thread, int)
	CheckThreadForUniq(ctx context.Context, thread models.Thread) ([]models.Thread, int)
	GetThreads(ctx context.Context, slug string, params models.RequestParameters) ([]models.Thread, error)
	GetThreadBySlug(ctx context.Context, slug string) (models.Thread, error)
	GetThreadById(ctx context.Context, id int) (models.Thread, error)
	CreatePosts(ctx context.Context, posts []models.Post, thread models.Thread) ([]models.Post, int)
	ChangeVote(ctx context.Context, vote models.Vote, thread models.Thread) (models.Thread, error)
	ChangeThreadInfo(ctx context.Context, thread models.Thread) (models.Thread, int)
	GetUsers(ctx context.Context, slug string, params models.RequestParameters) ([]models.User, error)
}

type UseCase interface {
	CreateUser(ctx context.Context, user models.User) ([]models.User, error)
	GetUser(ctx context.Context, nickname string) (models.User, error)
	ChangeUserInfo(ctx context.Context, user models.User) (models.User, int)
	CreateForum(ctx context.Context, forum models.Forum) ([]models.Forum, int)
	GetForumDetails(ctx context.Context, slug string) (models.Forum, error)
	CreateThread(ctx context.Context, thread models.Thread) ([]models.Thread, int)
	GetThreads(ctx context.Context, slug string, params models.RequestParameters) ([]models.Thread, error)
	GetThreadBySlugOrId(ctx context.Context, slugOrId string) (models.Thread, error)
	CreatePosts(ctx context.Context, posts []models.Post, thread models.Thread) ([]models.Post, int)
	ChangeVote(ctx context.Context, vote models.Vote, thread models.Thread) (models.Thread, error)
	ChangeThreadInfo(ctx context.Context, newThread models.Thread, oldThread models.Thread) (models.Thread, int)
	GetUsers(ctx context.Context, slug string, params models.RequestParameters) ([]models.User, error)
}
