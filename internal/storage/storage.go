package storage

import "errors"

var (
	ErrSongNotFound       = errors.New("musik not found")
	ErrSongDetailNotFound = errors.New("musik detail not found")
	ErrSongAlreadyExists  = errors.New("musik exists")
)
