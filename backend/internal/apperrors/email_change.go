package apperrors

import (
	"backend/internal/apicodes"
	"github.com/pkg/errors"
)

type EmailChangeEmailAlreadyUsed struct {
	AppError
}
type EmailChangeSameEmailError struct {
	AppError
}

func (e *EmailChangeEmailAlreadyUsed) Error() string {
	return e.Description
}

func NewEmailChangeEmailAlreadyUsedError(desc string) *EmailChangeEmailAlreadyUsed {
	return &EmailChangeEmailAlreadyUsed{
		AppError: AppError{
			Code:        apicodes.API_Email_Change_Email_Already_Used,
			Description: desc,
		},
	}
}
func NewEmailChangeSameEmailError(desc string) *EmailChangeSameEmailError {
	return &EmailChangeSameEmailError{
		AppError: AppError{
			Code:        apicodes.API_Email_Change_Same_Email,
			Description: desc,
		},
	}
}

func IsNewEmailChangeEmailAlreadyUsedError(err error) bool {
	var takenErr *EmailChangeEmailAlreadyUsed
	return errors.As(err, &takenErr)
}

func IsNewEmailChangeSameEmailError(err error) bool {
	var sameErr *EmailChangeSameEmailError
	return errors.As(err, &sameErr)
}