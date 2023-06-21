package usecase

import (
	"context"
	"fmt"
	"github.com/BigBullas/TP_DB_project/internal/models"
	"github.com/BigBullas/TP_DB_project/internal/pkg/forume"
	"net/http"
	"strconv"
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

	if user.Email == "" {
		user.Email = thisUser.Email
	}
	if user.About == "" {
		user.About = thisUser.About
	}
	if user.FullName == "" {
		user.FullName = thisUser.FullName
	}

	usersWithSameInfo, _ := u.repo.CheckUserForUniq(ctx, user)
	if len(usersWithSameInfo) > 1 {
		return models.User{}, http.StatusConflict
	}
	return u.repo.ChangeUserInfo(ctx, user)
}

func (u *UseCase) CreateForum(ctx context.Context, forum models.Forum) ([]models.Forum, int) {
	forumsWithSameSlug, _ := u.repo.CheckForumForUniq(ctx, forum)
	if len(forumsWithSameSlug) > 0 {
		return forumsWithSameSlug, http.StatusConflict
	}

	author, err := u.repo.GetUser(ctx, forum.User)
	if err != nil {
		return []models.Forum{}, http.StatusInternalServerError
	}
	if author == (models.User{}) {
		return []models.Forum{}, http.StatusNotFound
	}
	forum.User = author.NickName
	return u.repo.CreateForum(ctx, forum)
}

func (u *UseCase) GetForumDetails(ctx context.Context, slug string) (models.Forum, error) {
	return u.repo.GetForumDetails(ctx, slug)
}

func (u *UseCase) CreateThread(ctx context.Context, thread models.Thread) ([]models.Thread, int) {
	if thread.Slug != "" {
		threadsWithSameSlug, _ := u.repo.CheckThreadForUniq(ctx, thread)
		if len(threadsWithSameSlug) > 0 {
			return threadsWithSameSlug, http.StatusConflict
		}
	}

	author, err := u.repo.GetUser(ctx, thread.Author)
	if err != nil {
		return []models.Thread{}, http.StatusInternalServerError
	}
	if author == (models.User{}) {
		return []models.Thread{}, http.StatusNotFound
	}
	thread.Author = author.NickName

	forum, err := u.repo.GetForumDetails(ctx, thread.Forum)
	if err != nil {
		return []models.Thread{}, http.StatusInternalServerError
	}
	if forum == (models.Forum{}) {
		return []models.Thread{}, http.StatusNotFound
	}
	thread.Forum = forum.Slug

	return u.repo.CreateThread(ctx, thread)
}

func (u *UseCase) GetThreads(ctx context.Context, slug string, params models.RequestParameters) ([]models.Thread, error) {
	thisForum, err := u.repo.GetForumDetails(ctx, slug)
	if err != nil {
		return []models.Thread{}, models.InternalError
	}
	if thisForum == (models.Forum{}) {
		return []models.Thread{}, models.NotFound
	}
	return u.repo.GetThreads(ctx, slug, params)
}

func (u *UseCase) GetThreadBySlugOrId(ctx context.Context, slugOrId string) (models.Thread, error) {
	var thisThread models.Thread
	var errSlug error
	var errId error

	slugOrIdNum, err := strconv.Atoi(slugOrId)
	if err != nil {
		thisThread, errSlug = u.repo.GetThreadBySlug(ctx, slugOrId)
	} else {
		thisThread, errId = u.repo.GetThreadById(ctx, slugOrIdNum)
	}
	if errSlug != nil || errId != nil {
		fmt.Println("useCase slugOrId ", errSlug, errId)
		return models.Thread{}, models.InternalError
	}
	if thisThread == (models.Thread{}) {
		fmt.Println("delivery not fount thisThread")
		return models.Thread{}, models.NotFound
	}
	return thisThread, nil
}

func (u *UseCase) CreatePosts(ctx context.Context, posts []models.Post, thread models.Thread) ([]models.Post, int) {
	return u.repo.CreatePosts(ctx, posts, thread)
}
