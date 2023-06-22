package models

// easyjson -all ./internal/models/post_detailed.go

type PostDetailed struct {
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
	Author *User   `json:"author,omitempty"`
	Post   Post    `json:"post,omitempty"`
}
