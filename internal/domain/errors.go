package domain

import "errors"

var ErrEmailAlreadyExists = errors.New("email already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrNotFound = errors.New("not found")
