package apierrors

import (
	"errors"
)

var ErrUnrecoverableError = errors.New("fatal unrecoverable error")

func IsUnrecoverableError(err error) bool {
	return errors.Is(err, ErrUnrecoverableError)
}
