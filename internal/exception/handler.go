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
	ErrBadRequest          = &ApiError{Status: 400, Code: "BAD_REQUEST", Message: "invalid request"}
	ErrBadRequestID        = &ApiError{Status: 400, Code: "BAD_REQUEST", Message: "id is required"}
	ErrParseUUIDFailed     = &ApiError{Status: 400, Code: "BAD_REQUEST", Message: "UUID format is invalid"}
	ErrInvalidPassword     = &ApiError{Status: 400, Code: "BAD_REQUEST", Message: "password is invalid"}
	ErrUnauthorized        = &ApiError{Status: 401, Code: "UNAUTHORIZED", Message: "authentication required"}
	ErrInvalidBearerToken  = &ApiError{Status: 401, Code: "UNAUTHORIZED", Message: "bearer Token format is invalid"}
	ErrExpiredBearerToken  = &ApiError{Status: 401, Code: "UNAUTHORIZED", Message: "token expired or invalid"}
	ErrInvalidPayloadToken = &ApiError{Status: 401, Code: "UNAUTHORIZED", Message: "payload token is invalid"}
)

var (
	ErrSysExpiredOrInvalidBearerToken = fmt.Errorf("token expired or invalid")
	ErrSysTokenAssigningInvalid       = fmt.Errorf("token assigning is invalid")
	ErrSysAssigningToken              = fmt.Errorf("error on assigning token")
	ErrSysInvalidPayloadToken         = fmt.Errorf("payload token is invalid")
	ErrSysCustomerNotFound            = fmt.Errorf("customer not found")
	ErrSysInvalidPassword             = fmt.Errorf("password is invalid")
	ErrSysHashPassword                = fmt.Errorf("error on hash password")
	ErrSysSaveRefreshToken            = fmt.Errorf("error on save refresh token")
	ErrSysGeneratePairTokens          = fmt.Errorf("error on generate pair tokens")
	ErrSysGenerateRandomDates         = fmt.Errorf("error on generate random dates")
	ErrSysParseTimeDate               = fmt.Errorf("error on parse date to time.Time")
	ErrSysSaveCustomer                = fmt.Errorf("error on save customer in database")
	ErrSysLoadCustomerList            = fmt.Errorf("error on load customer list in database")
	ErrSysFindCustomerById            = fmt.Errorf("error on load customer by id in database")
	ErrSysFindCustomerByEmail         = fmt.Errorf("error on load customer by email in database")
	ErrSysDeleteCustomer              = fmt.Errorf("error on delete customer by id in database")
	ErrSysUpdateCustomer              = fmt.Errorf("error on update customer by id in database")
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
