package delivery

import (
	"fmt"
	"github.com/BigBullas/TP_DB_project/internal/models"
	User "github.com/BigBullas/TP_DB_project/internal/pkg/forume"
	"github.com/BigBullas/TP_DB_project/internal/utils"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"net/http"
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
		utils.Response(w, http.StatusNotFound, nil)
		return
	}

	user := models.User{}
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil) // почему здесь StatusInternalServerError
		return
	}
	user.NickName = nickname

	finalUser, err := h.uc.CreateUser(r.Context(), user)
	if err == nil {
		newUser := finalUser[0]
		utils.Response(w, http.StatusCreated, newUser)
		return
	}
	utils.Response(w, http.StatusConflict, finalUser)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname, flag := vars["nickname"]
	if !flag {
		utils.Response(w, http.StatusNotFound, nil)
		return
	}

	foundUser, err := h.uc.GetUser(r.Context(), nickname)
	if err == nil && foundUser == (models.User{}) {
		utils.Response(w, http.StatusNotFound, nickname)
		return
	}
	if err == nil {
		utils.Response(w, http.StatusOK, foundUser)
		return
	}
	utils.Response(w, http.StatusNotFound, nickname)
}

func (h *Handler) ChangeUserInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname, flag := vars["nickname"]
	if !flag {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	user := models.User{}
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil) // почему здесь StatusInternalServerError
		return
	}
	user.NickName = nickname
	fmt.Println("delivery after unmarshal user:", user)
	changedUser, status := h.uc.ChangeUserInfo(r.Context(), user)
	fmt.Println("delivery after useCase user:", changedUser)
	utils.Response(w, status, changedUser)
}

func (h *Handler) CreateForum(w http.ResponseWriter, r *http.Request) {
	forum := models.Forum{}
	err := easyjson.UnmarshalFromReader(r.Body, &forum)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	createdForums, status := h.uc.CreateForum(r.Context(), forum)
	if len(createdForums) > 0 {
		utils.Response(w, status, createdForums[0])
		return
	}
	utils.Response(w, status, forum.Title)
}
