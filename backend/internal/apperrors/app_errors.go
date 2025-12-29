package apperrors

import (
	"backend/internal/apicodes"
	"github.com/pkg/errors"
)

type AppError struct {
	Code        int
	Description string
}

type AppInvalidInputError struct {
	AppError
	Field string
}

func (e *AppError) Error() string {
	return e.Description
}

func (e *AppInvalidInputError) Error() string {
	return e.Description
}

func NewInvalidInputError(field, description string) *AppInvalidInputError {
	return &AppInvalidInputError{
		AppError: AppError{Code: apicodes.API_General_Invalid_Input_Value, Description: description},
		Field:    field,
	}
}
func NewGeneralCustomError(description string) *AppError {
	return &AppError{
		Code: apicodes.API_General_Custom_Error,
		Description: description,
	}
}

func NewAppError(code int, description string) *AppError {
	return &AppError{
		Code:        code,
		Description: description,
	}
}

func IsAppInvalidInputError(err error) bool {
	var invalidInputErr *AppInvalidInputError
	return errors.As(err, &invalidInputErr)
}
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}