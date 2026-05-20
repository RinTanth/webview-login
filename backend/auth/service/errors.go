package service

import "errors"

var (
	ErrPasswordMismatch   = errors.New("passwords do not match")
	ErrConflict           = errors.New("username or email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
