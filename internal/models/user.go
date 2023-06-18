package models

// easyjson -all ./internal/models/forume.go

type User struct {
	ID       int    `json:"-"`
	NickName string `json:"nickname,omitempty"`
	FullName string `json:"fullname"`
	About    string `json:"about,omitempty"`
	Email    string `json:"email"`
}
