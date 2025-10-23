package util

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func FormatValidationError(err error) []ValidationError {
	var ve validator.ValidationErrors
	var validationErrors []ValidationError

	if errors.As(err, &ve) {
		for _, fieldError := range ve {
			var errorMessage string
			switch fieldError.Tag() {
			case "required":
				errorMessage = "can't be null"
			case "email":
				errorMessage = "must be a valid email"
			case "min":
				errorMessage = "must be at least " + fieldError.Param() + " characters long"
			case "max":
				errorMessage = "must be at most " + fieldError.Param() + " characters long"
			default:
				errorMessage = "invalid value"
			}
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldError.Field(),
				Error: errorMessage,
			})

		}
	}
	return validationErrors
}
