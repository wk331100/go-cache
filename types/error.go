package types

import "errors"

var (
	ErrKeyNotExist = errors.New("key not exist")
	ErrEmptyList   = errors.New("list is empty")
	ErrStartStop   = errors.New("start or stop is invalid")
	ErrHashKey     = errors.New("hash key is not exist")
	ErrHashField   = errors.New("hash field is not exist")
	ErrSetKey      = errors.New("set key is not exist")
	ErrZSetKey     = errors.New("zset key is not exist")
)
