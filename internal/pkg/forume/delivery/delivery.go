package delivery

import (
	"encoding/json"
	"github.com/BigBullas/TP_DB_project/internal/models"
	User "github.com/BigBullas/TP_DB_project/internal/pkg/forume"
	"github.com/BigBullas/TP_DB_project/internal/utils"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"net/http"
	"strconv"
)

type Handler struct {
	uc User.UseCase
}

func NewForumHandler(useCase User.UseCase) *Handler {
	return &Handler{uc: useCase}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname, flag := vars["nickname"]
	if !flag {
		utils.Response(w, http.StatusNotFound, nil, false)
		return
	}

	user := models.User{}
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil, false) // почему здесь StatusInternalServerError
		return
	}
	user.NickName = nickname

	finalUser, err := h.uc.CreateUser(r.Context(), user)
	if err == nil {
		newUser := finalUser[0]
		utils.Response(w, http.StatusCreated, newUser, false)
		return
	}
	utils.Response(w, http.StatusConflict, finalUser, false)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname, flag := vars["nickname"]
	if !flag {
		utils.Response(w, http.StatusNotFound, nil, false)
		return
	}

	foundUser, err := h.uc.GetUser(r.Context(), nickname)
	if err == nil && foundUser == (models.User{}) {
		utils.Response(w, http.StatusNotFound, nickname, false)
		return
	}
	if err == nil {
		utils.Response(w, http.StatusOK, foundUser, false)
		return
	}
	utils.Response(w, http.StatusNotFound, nickname, false)
}

func (h *Handler) ChangeUserInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname, flag := vars["nickname"]
	if !flag {
		utils.Response(w, http.StatusBadRequest, nil, false)
		return
	}

	user := models.User{}
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil, false) // почему здесь StatusInternalServerError
		return
	}
	user.NickName = nickname
	changedUser, status := h.uc.ChangeUserInfo(r.Context(), user)
	utils.Response(w, status, changedUser, false)
}

func (h *Handler) CreateForum(w http.ResponseWriter, r *http.Request) {
	forum := models.Forum{}
	err := easyjson.UnmarshalFromReader(r.Body, &forum)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil, false)
		return
	}

	createdForums, status := h.uc.CreateForum(r.Context(), forum)
	if len(createdForums) > 0 {
		utils.Response(w, status, createdForums[0], false)
		return
	}
	utils.Response(w, status, forum.Title, false)
}

func (h *Handler) GetForumDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, flag := vars["slug"]
	if !flag {
		utils.Response(w, http.StatusBadRequest, nil, false)
		return
	}

	foundForum, err := h.uc.GetForumDetails(r.Context(), slug)
	if err == nil && foundForum == (models.Forum{}) {
		utils.Response(w, http.StatusNotFound, slug, false)
		return
	}
	if err == nil {
		utils.Response(w, http.StatusOK, foundForum, false)
		return
	}
	utils.Response(w, http.StatusNotFound, slug, false)
}

func (h *Handler) CreateThread(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, flag := vars["slug"]
	if !flag {
		utils.Response(w, http.StatusBadRequest, nil, false)
		return
	}

	thread := models.Thread{}
	err := easyjson.UnmarshalFromReader(r.Body, &thread)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil, false)
		return
	}
	thread.Forum = slug

	createdThreads, status := h.uc.CreateThread(r.Context(), thread)
	if len(createdThreads) > 0 {
		utils.Response(w, status, createdThreads[0], false)
		return
	}
	utils.Response(w, status, thread.Title, false)
}

func (h *Handler) GetThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, flag := vars["slug"]
	if !flag {
		utils.Response(w, http.StatusBadRequest, nil, false)
		return
	}
	limitInput := r.URL.Query().Get("limit")
	sinceInput := r.URL.Query().Get("since")
	descInput := r.URL.Query().Get("desc")

	params := models.RequestParameters{}
	if limitInput == "" {
		params.Limit = 100
	} else {
		limit, errLimit := strconv.Atoi(limitInput)
		if errLimit != nil {
			utils.Response(w, http.StatusNotFound, slug, false)
			return
		}
		params.Limit = limit
	}

	params.Since = sinceInput

	if descInput == "" {
		params.Desc = false
	} else {
		desc, errDesc := strconv.ParseBool(descInput)
		if errDesc != nil {
			utils.Response(w, http.StatusNotFound, slug, false)
			return
		}
		params.Desc = desc
	}

	foundThreads, err := h.uc.GetThreads(r.Context(), slug, params)
	if err == models.NotFound {
		utils.Response(w, http.StatusNotFound, slug, false)
		return
	}
	if err == nil && len(foundThreads) == 0 {
		utils.Response(w, http.StatusOK, []models.Thread{}, false)
		return
	}
	if err == nil {
		utils.Response(w, http.StatusOK, foundThreads, false)
		return
	}
	utils.Response(w, http.StatusNotFound, slug, false)
}

func (h *Handler) CreatePosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId, flag := vars["slug_or_id"]
	if !flag {
		utils.Response(w, http.StatusBadRequest, nil, false)
		return
	}

	thisThread, errThread := h.uc.GetThreadBySlugOrId(r.Context(), slugOrId)
	if errThread == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil, false)
		return
	}
	if errThread == models.NotFound {
		utils.Response(w, http.StatusNotFound, slugOrId, false)
		return
	}

	var posts []models.Post
	decoder := json.NewDecoder(r.Body)
	errDec := decoder.Decode(&posts)
	if errDec != nil {
		utils.Response(w, http.StatusBadRequest, nil, false)
		return
	}

	if len(posts) == 0 {
		utils.Response(w, http.StatusCreated, []models.Post{}, false)
		return
	}

	createdPosts, status := h.uc.CreatePosts(r.Context(), posts, thisThread)
	if status == http.StatusConflict {
		utils.Response(w, status, slugOrId, true)
		return
	}
	utils.Response(w, status, createdPosts, false)
}

func (h *Handler) ChangeVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId, flag := vars["slug_or_id"]
	if !flag {
		utils.Response(w, http.StatusBadRequest, nil, false)
		return
	}

	thisThread, errThread := h.uc.GetThreadBySlugOrId(r.Context(), slugOrId)
	if errThread == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil, false)
		return
	}
	if errThread == models.NotFound {
		utils.Response(w, http.StatusNotFound, slugOrId, false)
		return
	}

	vote := models.Vote{}
	err := easyjson.UnmarshalFromReader(r.Body, &vote)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil, false)
		return
	}
	vote.Thread = thisThread.ID

	thisUser, err := h.uc.GetUser(r.Context(), vote.Nickname)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil, false)
		return
	}
	if thisUser == (models.User{}) {
		utils.Response(w, http.StatusNotFound, thisUser, false)
		return
	}

	changedThread, err := h.uc.ChangeVote(r.Context(), vote, thisThread)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil, false)
		return
	}
	if changedThread != (models.Thread{}) {
		utils.Response(w, http.StatusOK, changedThread, false)
		return
	}
	finalThread, errThread := h.uc.GetThreadBySlugOrId(r.Context(), slugOrId)
	if errThread == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil, false)
		return
	}
	utils.Response(w, http.StatusOK, finalThread, false)
}

func (h *Handler) GetThreadDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId, flag := vars["slug_or_id"]
	if !flag {
		utils.Response(w, http.StatusBadRequest, nil, false)
		return
	}

	foundThread, errThread := h.uc.GetThreadBySlugOrId(r.Context(), slugOrId)
	if errThread == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil, false)
		return
	}
	if errThread == models.NotFound {
		utils.Response(w, http.StatusNotFound, slugOrId, false)
		return
	}
	utils.Response(w, http.StatusOK, foundThread, false)
}

func (h *Handler) ChangeThreadInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slugOrId, flag := vars["slug_or_id"]
	if !flag {
		utils.Response(w, http.StatusBadRequest, nil, false)
		return
	}

	thread := models.Thread{}
	err := easyjson.UnmarshalFromReader(r.Body, &thread)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil, false)
		return
	}

	foundThread, errThread := h.uc.GetThreadBySlugOrId(r.Context(), slugOrId)
	if errThread == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil, false)
		return
	}
	if errThread == models.NotFound {
		utils.Response(w, http.StatusNotFound, slugOrId, false)
		return
	}

	changedThread, status := h.uc.ChangeThreadInfo(r.Context(), thread, foundThread)
	utils.Response(w, status, changedThread, false)

}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, flag := vars["slug"]
	if !flag {
		utils.Response(w, http.StatusBadRequest, nil, false)
		return
	}
	limitInput := r.URL.Query().Get("limit")
	sinceInput := r.URL.Query().Get("since")
	descInput := r.URL.Query().Get("desc")

	params := models.RequestParameters{}
	if limitInput == "" {
		params.Limit = 100
	} else {
		limit, errLimit := strconv.Atoi(limitInput)
		if errLimit != nil {
			utils.Response(w, http.StatusNotFound, slug, false)
			return
		}
		params.Limit = limit
	}

	params.Since = sinceInput

	if descInput == "" {
		params.Desc = false
	} else {
		desc, errDesc := strconv.ParseBool(descInput)
		if errDesc != nil {
			utils.Response(w, http.StatusNotFound, slug, false)
			return
		}
		params.Desc = desc
	}

	foundUsers, err := h.uc.GetUsers(r.Context(), slug, params)
	if err == models.NotFound {
		utils.Response(w, http.StatusNotFound, slug, false)
		return
	}
	if err == nil && len(foundUsers) == 0 {
		utils.Response(w, http.StatusOK, []models.User{}, false)
		return
	}
	if err == nil {
		utils.Response(w, http.StatusOK, foundUsers, false)
		return
	}
	utils.Response(w, http.StatusNotFound, slug, false)
}
