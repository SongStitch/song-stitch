package session

import "errors"

var ErrTooManyImages = errors.New("too many images requested")

var ErrUserNotFound = errors.New("user not found")
