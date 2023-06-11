package constants

import "errors"

var ErrTooManyImages = errors.New("too many images requested")

var ErrUserNotFound = errors.New("user not found")

var ErrNoImageFound = errors.New("no image found")
