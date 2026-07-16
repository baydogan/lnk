package errs

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrAlreadyExists     = errors.New("already exists")
	ErrInvalidURL        = errors.New("invalid url")
	ErrAliasExists       = errors.New("alias already taken")
	ErrURLLimit          = errors.New("url limit reached")
	ErrExpireFormat      = errors.New("invalid expire format (use: 1h, 7d, 30d)")
	ErrConfigNotFound    = errors.New("config not found")
	ErrNotLoggedIn       = errors.New("not logged in")
	ErrInvalidKey        = errors.New("invalid key")
	ErrInvalidUsername   = errors.New("invalid username")
	ErrInvalidRole       = errors.New("invalid role (use: admin, user)")
	ErrCannotDeleteAdmin = errors.New("cannot delete an admin user")
)
