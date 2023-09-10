package errorwrap

import (
	"errors"
	"fmt"
)

type QueryError struct {
	Package    string
	ResponeErr string
	err        error
	WrapError  error
}

var ErrNotFound = errors.New("not found")

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
