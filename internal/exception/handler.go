package exception

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ApiError) Error() string {
	return e.Message
}

// HTTP Mapping Errors
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
)

// System Mapping Errors
var (
	ErrSysNotFound          = errors.New("resource not found")
	ErrSysAlreadyExists     = errors.New("resource already exists")
	ErrSysInvalidInput      = errors.New("invalid input")
	ErrSysInvalidCredential = errors.New("invalid credential")
	ErrSysUnauthorized      = errors.New("unauthorized")
	ErrSysForbidden         = errors.New("forbidden")

	ErrSysInvalidToken  = errors.New("invalid token")
	ErrSysExpiredToken  = errors.New("expired token")
	ErrSysTokenNotFound = errors.New("token not found")

	ErrSysHash     = errors.New("hash operation failed")
	ErrSysGenerate = errors.New("generation failed")
	ErrSysParse    = errors.New("parse failed")

	ErrSysDatabase    = errors.New("database operation failed")
	ErrSysPersistence = errors.New("persistence operation failed")

	ErrSysInternal = errors.New("internal server error")
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var appErr *ApiError
		if errors.As(err, &appErr) {
			c.JSON(appErr.Status, gin.H{"code": appErr.Code, "error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL_ERROR", "message": "an unexpected error occurred"})
		}
	}
}

func ErrorWithContext(c *gin.Context, err ApiError) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":  err.Code,
		"error": fmt.Sprintf("%v error: %s", err.Message, err.Message),
	})
}

func BadRequestErrorHandler(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":   "BAD_REQUEST",
		"errors": HandleValidationError(err),
	})
}
