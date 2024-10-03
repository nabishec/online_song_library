package storage

import "errors"

var (
	ErrMusikNotFound       = errors.New("musik not found")
	ErrMusikDetailNotFound = errors.New("musik detail not found")
	ErrMusikAlreadyExists  = errors.New("musik exists")
)
