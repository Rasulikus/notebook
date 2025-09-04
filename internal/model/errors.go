package model

import "errors"

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrGenerateRefreshToken = errors.New("failed to generate refresh token after retries")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrWrongCredetials      = errors.New("wrong email or password")
)
