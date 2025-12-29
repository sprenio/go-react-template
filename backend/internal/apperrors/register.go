package apperrors

import (
	"backend/internal/apicodes"

	"github.com/pkg/errors"
)

type RegisterUserNameOrEmailTakenError struct {
	AppError
}

func (e *RegisterUserNameOrEmailTakenError) Error() string {
	return e.Description
}

func NewRegisterUserNameOrEmailTakenError(desc string) *RegisterUserNameOrEmailTakenError {
	return &RegisterUserNameOrEmailTakenError{
		AppError: AppError{
			Code:        apicodes.API_Register_User_Name_Or_Email_Taken,
			Description: desc,
		},
	}
}

func IsRegisterUserNameOrEmailTakenError(err error) bool {
	var takenErr *RegisterUserNameOrEmailTakenError
	return errors.As(err, &takenErr)
}
