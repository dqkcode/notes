package user

import (
	"errors"
)

var (
	ErrUserNotFound = errors.New("User not found")
	ErrUserAlreadyExist= errors.New("User already exist")
)