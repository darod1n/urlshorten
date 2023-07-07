package errstorage

import "errors"

var (
	ExistURL  = errors.New("url exist")
	RemoveURL = errors.New("url removed")
	NYI       = errors.New("not yet implemented")
)
