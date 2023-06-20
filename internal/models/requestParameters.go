package models

// easyjson -all ./internal/models/requestParameters.go

type RequestParameters struct {
	Desc  bool   `json:"desc"`
	Limit int    `json:"limit"`
	Since string `json:"since"`
}
