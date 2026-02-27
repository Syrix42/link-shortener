package auth

import "errors"

var ErrUserAlreadyExists = errors.New("Error user Already Exists")
var ErrInvalidEmailFormat = errors.New("Email format is Invalid")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidPassword = errors.New("Password is Invalid")
var ErrTooManyActiveSessions = errors.New("active sessions cannot be more than 5")
