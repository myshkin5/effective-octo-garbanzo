package services

import (
	"strings"
)

type ValidationError struct {
	errors []string
}

func NewValidationError(errors ...string) ValidationError {
	return ValidationError{
		errors: errors,
	}
}

func (e ValidationError) Error() string {
	return strings.Join(e.errors, ", ")
}

func (e ValidationError) Errors() []string {
	return e.errors
}
