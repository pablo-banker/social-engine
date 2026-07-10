package apiErrors

import "fmt"

func New(code int, message string, httpStatus int) *APIError {
	return &APIError{
		Code:       formatCode(code),
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

func formatCode(code int) string {
	return fmt.Sprintf("SE-%06d", code)
}
