package apperrors

import "errors"

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
