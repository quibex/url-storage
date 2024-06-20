package cache

import "errors"

var (
	ErrAlreadyExist = errors.New("alias already exist")
	ErrNotFound     = errors.New("not found")
)
