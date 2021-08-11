package http

import "errors"

var (
	ErrNotAuthenticated = errors.New("Not authenticated")
	ErrMethodNotAlowed  = errors.New("method not allowed")
	ErrDataIsNotJson    = errors.New("data not json")
)
