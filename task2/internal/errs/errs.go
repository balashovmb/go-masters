package errs

import (
	"fmt"
)

type ErrBadRequest struct {
	massage string
}

type ErrNotFound struct {
	Pair string
}

type ErrInternalServerError struct{}

func (err ErrBadRequest) Error() string {
	return err.massage
}

func (err ErrNotFound) Error() string {
	return fmt.Sprintf("pair %s not found", err.Pair)
}

func (err ErrInternalServerError) Error() string {
	return "internal server error"
}

func NewErrNotFound(pair string) error {
	return ErrNotFound{Pair: pair}
}

func NewErrBadRequest(massage string) error {
	return ErrBadRequest{massage: massage}
}
