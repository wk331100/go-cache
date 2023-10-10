package go_cache

import "errors"

var (
	ErrKeyNotExist = errors.New("key not exist")
	ErrEmptyList   = errors.New("list is empty")
	ErrStartLen    = errors.New("start or len is invalid")
	ErrHashKey     = errors.New("hash key is not exist")
	ErrHashField   = errors.New("hash field is not exist")
)
