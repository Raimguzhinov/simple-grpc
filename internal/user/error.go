package user

import "errors"

var (
	ErrBadFormat    = errors.New("bad promt format")
	ErrBadDateTime  = errors.New("bad datetime format")
	ErrUnexpected   = errors.New("unexpected error")
	ErrNotFound     = errors.New("not found")
	ErrBadProcedure = errors.New("bad procedure name")
)
