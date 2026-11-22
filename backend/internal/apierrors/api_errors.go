package apierrors

import (
	"backend/internal/apicodes"
	"errors"
	"reflect"
	"strings"
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

func NewInvalidInputError[T any](field, description string, requestData *T) *AppInvalidInputError {
	fieldName := field

	if requestData != nil {
		// reflect.TypeOf(*requestData) pobiera typ T (struktury), a nie wskaźnika *T
		t := reflect.TypeOf(*requestData)

		// Obsługa sytuacji, gdyby T samo w sobie było interfejsem lub wskaźnikiem
		if t != nil && t.Kind() == reflect.Struct {
			if structField, found := t.FieldByName(field); found {
				jsonTag := structField.Tag.Get("json")

				// Jeśli tag istnieje i nie jest ignorowany (`json:"-"`)
				if jsonTag != "" && jsonTag != "-" {
					// Odrzucamy opcje takie jak omitempty
					name := strings.Split(jsonTag, ",")[0]
					if name != "" {
						fieldName = name
					}
				}
			}
		}
	}

	return &AppInvalidInputError{
		AppError: AppError{
			Code:        apicodes.API_General_Invalid_Input_Value,
			Description: description,
		},
		Field: fieldName,
	}
}
func NewGeneralCustomError(description string) *AppError {
	return &AppError{
		Code:        apicodes.API_General_Custom_Error,
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
