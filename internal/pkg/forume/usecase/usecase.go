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
		fmt.Println("useCase get thread error ", errSlug, errId, thisThread)
		return models.Thread{}, models.InternalError
	}
	if thisThread == (models.Thread{}) {
		return models.Thread{}, models.NotFound
	}
	return thisThread, nil
}

func (u *UseCase) CreatePosts(ctx context.Context, posts []models.Post, thread models.Thread) ([]models.Post, int) {
	return u.repo.CreatePosts(ctx, posts, thread)
}

func (u *UseCase) ChangeVote(ctx context.Context, vote models.Vote, thread models.Thread) (models.Thread, error) {
	return u.repo.ChangeVote(ctx, vote, thread)
}

func (u *UseCase) ChangeThreadInfo(ctx context.Context, newThread models.Thread, oldThread models.Thread) (models.Thread, int) {
	changeFlag := false
	if newThread.Title != "" {
		changeFlag = true
		oldThread.Title = newThread.Title
	}
	if newThread.Message != "" {
		changeFlag = true
		oldThread.Message = newThread.Message
	}
	if !changeFlag {
		return oldThread, http.StatusOK
	}
	return u.repo.ChangeThreadInfo(ctx, oldThread)
}

func (u *UseCase) GetUsers(ctx context.Context, slug string, params models.RequestParameters) ([]models.User, error) {
	thisForum, err := u.repo.GetForumDetails(ctx, slug)
	if err != nil {
		return []models.User{}, models.InternalError
	}
	if thisForum == (models.Forum{}) {
		return []models.User{}, models.NotFound
	}
	return u.repo.GetUsers(ctx, slug, params)
}

func (u *UseCase) GetPostDetails(ctx context.Context, id int, related []string) (models.PostDetailed, error) {
	return u.repo.GetPostDetails(ctx, id, related)
}

func (u *UseCase) ChangePostInfo(ctx context.Context, newPost models.Post, oldPost models.Post) (models.Post, int) {
	if newPost.Message == "" {
		return oldPost, http.StatusOK
	}
	if newPost.Message == oldPost.Message {
		return oldPost, http.StatusOK
	}
	oldPost.Message = newPost.Message
	oldPost.IsEdited = true
	return u.repo.ChangePostInfo(ctx, oldPost)
}

func (u *UseCase) GetStatus(ctx context.Context) (models.Info, int) {
	return u.repo.GetStatus(ctx)
}

func (u *UseCase) Clear(ctx context.Context) int {
	return u.repo.Clear(ctx)
}

func (u *UseCase) GetPosts(ctx context.Context, threadID int, params models.RequestParameters) ([]models.Post, error) {
	switch params.Sort {
	case "flat":
		return u.repo.GetPostsFlat(ctx, params, threadID)
	case "tree":
		return u.repo.GetPostsTree(ctx, params, threadID)
	case "parent_tree":
		return u.repo.GetPostsParent(ctx, params, threadID)
	default:
		return []models.Post{}, models.InternalError
	}
}
