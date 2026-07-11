package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var jwtRegex = regexp.MustCompile(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+$`)

func jwtValidator(fl validator.FieldLevel) bool {
	token := fl.Field().String()
	return jwtRegex.MatchString(token)
}
