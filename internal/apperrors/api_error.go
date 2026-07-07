package apperrors

type ApiError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func (e *ApiError) Error() string {
	return e.Message
}

func (e *ApiError) WithMessage(message string) *ApiError {
	clone := *e
	clone.Message = message
	return &clone
}

func (e *ApiError) WithDetails(details any) *ApiError {
	clone := *e
	clone.Details = details
	return &clone
}
