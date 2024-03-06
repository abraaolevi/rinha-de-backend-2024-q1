package internal

import "net/http"

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

var ErrAccountNotFound = NewError("account not found", http.StatusNotFound)
var ErrInsufficientBalance = NewError("insufficient account balance", http.StatusUnprocessableEntity)
