package apiErrors

type APIError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	Internal   string `json:"-"`
	HTTPStatus int    `json:"-"`
}

func (e *APIError) Error() string {
	return e.Message
}

func (e *APIError) WithInternal(internalMsg string) *APIError {
	return &APIError{
		Code:       e.Code,
		Message:    e.Message,
		Internal:   internalMsg,
		HTTPStatus: e.HTTPStatus,
		Details:    e.Details,
	}
}

func (e *APIError) WithDetails(details string) *APIError {
	return &APIError{
		Code:       e.Code,
		Message:    e.Message,
		Details:    details,
		Internal:   e.Internal,
		HTTPStatus: e.HTTPStatus,
	}
}

func (e *APIError) SetHTTPStatus(statusCode int) *APIError {
	return &APIError{
		Code:       e.Code,
		Message:    e.Message,
		Internal:   e.Internal,
		HTTPStatus: statusCode,
		Details:    e.Details,
	}
}
