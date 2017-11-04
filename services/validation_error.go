package services

type ValidationError struct {
	errors map[string][]string
}

func NewValidationError(errors map[string][]string) ValidationError {
	return ValidationError{
		errors: errors,
	}
}

func (e ValidationError) Error() string {
	message := "Validation error: "
	for key, values := range e.errors {
		for _, value := range values {
			message += key + " " + value + ", "
		}
	}
	return message
}

func (e ValidationError) Errors() map[string][]string {
	return e.errors
}
