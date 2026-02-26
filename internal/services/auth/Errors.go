package auth

import "errors"

var ErrUserAlreadyExists = errors.New("Error user Already Exists")
var ErrInvalidEmailFormat = errors.New("Email format is Invalid")
