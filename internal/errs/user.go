package errs

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrWrongCredetials   = errors.New("wrong email or password")
)
