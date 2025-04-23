package errs

import (
	"fmt"
)

type ErrBadRequest struct{}

type ErrNotFound struct {
	Pair string
}

type ErrInternalServerError struct{}

func (err ErrBadRequest) Error() string {
	return "bad request"
}

func (err ErrNotFound) Error() string {
	return fmt.Sprintf("pair %s not found", err.Pair)
}

func (err ErrInternalServerError) Error() string {
	return "internal server error"
}
