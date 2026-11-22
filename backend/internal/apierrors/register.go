package apierrors

import (
	"backend/internal/apicodes"

	"errors"
)

type RegisterUserNameTakenError struct {
	AppError
}

func (e *RegisterUserNameTakenError) Error() string {
	return e.Description
}

func NewRegisterUserNameTakenError(desc string) *RegisterUserNameTakenError {
	return &RegisterUserNameTakenError{
		AppError: AppError{
			Code:        apicodes.API_Register_User_Name_Taken,
			Description: desc,
		},
	}
}

func IsRegisterUserNameTakenError(err error) bool {
	var takenErr *RegisterUserNameTakenError
	return errors.As(err, &takenErr)
}
