package auth

import "errors"

var ErrTokenExpired = errors.New("token expired")
var ErrTokenInvalid = errors.New("token invalid")
