package middleware

import (
	"errors"
	"strings"

	"github.com/brunosilv96/bs-aesthetics-api/internal/auth"
	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
	"github.com/gin-gonic/gin"
)

type AccessTokenMiddleware struct {
	tokenService auth.TokenService
}

func NewAccessTokenMiddleware(tokenService auth.TokenService) *AccessTokenMiddleware {
	return &AccessTokenMiddleware{
		tokenService: tokenService,
	}
}

func (middleware *AccessTokenMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extracts the bearer token from the header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(exception.ErrUnauthorized)
			c.Abort()
			return
		}

		// 2. Verify if format is Bearer "<token>"
		const bearerPrefix = "Bearer "

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.Error(exception.ErrInvalidBearerToken)
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, bearerPrefix)

		// 3. Check and Validate Access Token
		claims, err := middleware.tokenService.ValidateAccessToken(token)
		if err != nil {
			if errors.Is(err, exception.ErrSysExpiredOrInvalidBearerToken) {
				c.Error(exception.ErrExpiredBearerToken)
				c.Abort()
				return
			}

			c.Error(exception.ErrInvalidBearerToken)
			c.Abort()
			return
		}

		c.Set("customer_id", claims.CustomerID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
