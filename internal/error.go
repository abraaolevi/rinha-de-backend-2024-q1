package internal

type Error struct {
	err    string
	Status int
}

func NewError(message string, status int) error {
	return &Error{
		err:    message,
		Status: status,
	}
}

func (e *Error) Error() string {
	return e.err
}
