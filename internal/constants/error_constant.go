package constants

import "errors"

var (
	ErrNotLoggedIn = errors.New("not logged in and no permission")
)