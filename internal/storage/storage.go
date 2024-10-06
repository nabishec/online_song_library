package storage

import "errors"

var (
	ErrSongNotFound       = errors.New("song not found")
	ErrSongDetailNotFound = errors.New("song detail not found")
	ErrSongAlreadyExists  = errors.New("song exists")
)
