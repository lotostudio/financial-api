package repo

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user doesn't exists")
	ErrSessionNotFound   = errors.New("session doesn't exists")

	ErrCurrencyNotFound = errors.New("currency doesn't exists")

	ErrAccountNotFound = errors.New("account doesn't exists")
)
