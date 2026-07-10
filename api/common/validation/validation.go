package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitValidator() {
	validate = validator.New()
	validate.RegisterValidation("jwt", jwtValidator)
	validate.RegisterValidation("avatarRef", avatarRefValidator)
}

func Struct(payload any) error {
	if err := validate.Struct(payload); err != nil {
		return parseError(err, payload)
	}

	return nil
}

func parseError(err error, input any) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			jsonField := getJSONName(input, e.StructField())

			switch e.Tag() {
			case "required":
				return fmt.Errorf("the field '%s' is required", jsonField)
			case "email":
				return fmt.Errorf("the field '%s' must be a valid email address", jsonField)
			case "min":
				return fmt.Errorf("the field '%s' must be at least %s characters long", jsonField, e.Param())
			case "max":
				return fmt.Errorf("the field '%s' must be at most %s characters long", jsonField, e.Param())
			case "len":
				return fmt.Errorf("the field '%s' must be exactly %s characters long", jsonField, e.Param())
			case "url":
				return fmt.Errorf("the field '%s' must be a valid URL", jsonField)
			case "jwt":
				return fmt.Errorf("the field '%s' must be a valid JWT token", jsonField)
			case "numeric":
				return fmt.Errorf("the field '%s' must contain only numbers", jsonField)
			default:
				return fmt.Errorf("the field '%s' is invalid", jsonField)
			}
		}
	}
	return fmt.Errorf("validation failed due to an unexpected error")
}
