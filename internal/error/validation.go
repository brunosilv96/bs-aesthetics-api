package error

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func HandleValidationError(err error) []ValidationError {
	validationErrors, ok := err.(validator.ValidationErrors)

	if !ok {
		return []ValidationError{
			{
				Field:   "request",
				Message: err.Error(),
			},
		}
	}

	errors := make([]ValidationError, 0, len(validationErrors))

	for _, fieldErr := range validationErrors {
		errors = append(errors, ValidationError{
			Field:   fieldErr.Field(),
			Message: getValidationMessage(fieldErr),
		})
	}

	return errors
}

func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "field is required"

	case "email":
		return "must be a valid email address"

	case "min":
		return fmt.Sprintf(
			"must be at least %s characters long",
			err.Param(),
		)

	case "max":
		return fmt.Sprintf(
			"must be at most %s characters long",
			err.Param(),
		)

	case "len":
		return fmt.Sprintf(
			"must be exactly %s characters long",
			err.Param(),
		)

	case "oneof":
		return fmt.Sprintf(
			"must be one of: %s",
			err.Param(),
		)

	case "uuid":
		return "must be a valid UUID"

	case "datetime":
		return "must be a valid date"

	case "gte":
		return fmt.Sprintf(
			"must be greater than or equal to %s",
			err.Param(),
		)

	case "lte":
		return fmt.Sprintf(
			"must be less than or equal to %s",
			err.Param(),
		)

	default:
		return "invalid value"
	}
}
