package repository

import "errors"

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrDuplicateUsernames = errors.New("this username is already taken")
	ErrDuplicateEmails    = errors.New("this email is already taken")
	ErrUnAuthorized       = errors.New("unauthorized")
)
