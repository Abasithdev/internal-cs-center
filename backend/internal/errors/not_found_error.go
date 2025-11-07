package errors

type NotFoundError struct {
	Msg string
}

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{Msg: msg}
}

func (notFoundErr *NotFoundError) Error() string {
	if notFoundErr.Msg != "" {
		return "Data not found" + notFoundErr.Msg
	}
	return "Data not found"
}
