package util

import "errors"

var (
	ErrorExpiredToken = errors.New("token has expired")
	ErrorMismatchSessionToken = errors.New("mismatch session token")
	ErrorBlockedSession = errors.New("blocked session")
	ErrorIncorrectUser = errors.New("incorrect user")
	ErrorExpiredSession = errors.New("expired session")
)