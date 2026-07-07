package apperrors

var (
	ErrBadRequest   = &ApiError{Status: 400, Code: "BAD_REQUEST", Message: "invalid request"}
	ErrInvalidInput = &ApiError{Status: 400, Code: "INVALID_INPUT", Message: "one or more fields are invalid"}
	ErrInvalidUUID  = &ApiError{Status: 400, Code: "INVALID_UUID", Message: "invalid UUID format"}

	ErrUnauthorized       = &ApiError{Status: 401, Code: "UNAUTHORIZED", Message: "authentication required"}
	ErrInvalidCredentials = &ApiError{Status: 401, Code: "INVALID_CREDENTIALS", Message: "invalid email or password"}
	ErrInvalidToken       = &ApiError{Status: 401, Code: "INVALID_TOKEN", Message: "token is invalid"}
	ErrExpiredToken       = &ApiError{Status: 401, Code: "TOKEN_EXPIRED", Message: "token has expired"}

	ErrForbidden = &ApiError{Status: 403, Code: "FORBIDDEN", Message: "access denied"}

	ErrNotFound = &ApiError{Status: 404, Code: "NOT_FOUND", Message: "resource not found"}

	ErrConflict = &ApiError{Status: 409, Code: "CONFLICT", Message: "resource already exists"}

	ErrUnprocessableEntity = &ApiError{
		Status:  422,
		Code:    "UNPROCESSABLE_ENTITY",
		Message: "request cannot be processed",
	}

	ErrInternal = &ApiError{
		Status:  500,
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "an unexpected error occurred",
	}
	ErrRequestTimeout = &ApiError{
		Status:  504,
		Code:    "REQUEST_TIMEOUT",
		Message: "request processing timed out",
	}
)
