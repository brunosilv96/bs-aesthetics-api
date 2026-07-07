package middleware

import (
	"strings"

	"github.com/brunosilv96/bs-aesthetics-api/internal/apperrors"
	"github.com/brunosilv96/bs-aesthetics-api/internal/model"
	"github.com/brunosilv96/bs-aesthetics-api/pkg"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenService pkg.TokenService
}

func NewAccessTokenMiddleware(tokenService pkg.TokenService) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
	}
}

func (middleware *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extracts the bearer token from the header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		// 2. Verify if format is Bearer "<token>"
		const bearerPrefix = "Bearer "

		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.Error(apperrors.ErrInvalidToken)
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, bearerPrefix)

		// 3. Check and Validate Access Token
		claims, err := middleware.tokenService.ValidateAccessToken(token)
		if err != nil {
			c.Error(apperrors.ErrExpiredToken)
			c.Abort()
			return
		}

		c.Set("auth_identity", &model.Identity{
			CustomerID: claims.CustomerID,
			Role:       claims.Role,
		})

		c.Next()
	}
}
