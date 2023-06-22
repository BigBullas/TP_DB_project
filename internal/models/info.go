package models

// easyjson -all ./internal/models/info.go

type Info struct {
	Users   int64 `json:"user"`
	Forums  int64 `json:"forum"`
	Threads int64 `json:"thread"`
	Posts   int64 `json:"post"`
}
