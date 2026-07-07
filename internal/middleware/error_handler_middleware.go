package middleware

import (
	"context"
	"errors"

	"github.com/brunosilv96/bs-aesthetics-api/internal/apperrors"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		for _, ginErr := range c.Errors {
			err := ginErr.Err

			switch {
			case errors.Is(err, context.Canceled):
				// Cliente cancelou a conexão.
				return

			case errors.Is(err, context.DeadlineExceeded):
				apperrors.WriteError(c, apperrors.ErrRequestTimeout)
				return
			}

			var apiErr *apperrors.ApiError
			if errors.As(err, &apiErr) {
				apperrors.WriteError(c, apiErr)
				return
			}
		}

		// Nenhum erro conhecido.
		apperrors.WriteError(c, apperrors.ErrInternal)
	}
}
