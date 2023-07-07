package errstorage

import "errors"

var (
	ErrExistURL  = errors.New("url exist")
	ErrRemoveURL = errors.New("url removed")
	ErrNYI       = errors.New("not yet implemented")
)
