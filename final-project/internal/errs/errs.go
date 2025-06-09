package errs

type ErrBadRequest struct {
	msg string
}

func (e ErrBadRequest) Error() string {
	return e.msg
}
