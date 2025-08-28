package errs

import "errors"

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrGenerateRefreshToken = errors.New("failed to generate refresh token after retries")
)
