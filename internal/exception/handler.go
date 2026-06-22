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

var (
	ErrNotFound            = &ApiError{Status: 404, Code: "NOT_FOUND", Message: "resource not found"}
	ErrCustomerNotFound    = &ApiError{Status: 404, Code: "NOT_FOUND", Message: "customer not found"}
	ErrUnauthorized        = &ApiError{Status: 401, Code: "UNAUTHORIZED", Message: "authentication required"}
	ErrBadRequest          = &ApiError{Status: 400, Code: "BAD_REQUEST", Message: "invalid request"}
	ErrBadRequestID        = &ApiError{Status: 400, Code: "BAD_REQUEST", Message: "id is required"}
	ErrParseUUIDFailed     = &ApiError{Status: 400, Code: "BAD_REQUEST", Message: "UUID format is invalid"}
	ErrInvalidPassword     = &ApiError{Status: 400, Code: "BAD_REQUEST", Message: "Password is invalid"}
	ErrSysCustomerNotFound = fmt.Errorf("Customer not found")
	ErrSysInvalidPassword  = fmt.Errorf("Password is invalid")
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

func BadRequestErrorHandler(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":   "BAD_REQUEST",
		"errors": HandleValidationError(err),
	})
}
