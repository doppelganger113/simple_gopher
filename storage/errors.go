package storage

import "errors"

var ErrDuplicate = errors.New("duplicate, already exists")

type NotFound struct {
	Msg string
}

func (nf NotFound) Error() string {
	return nf.Msg
}
