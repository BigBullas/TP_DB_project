package models

import "errors"

var (
	Conflict      = errors.New("conflict")
	InternalError = errors.New("InternalError")
	NotFoundForum = errors.New("NotFoundForum")
)
