package middleware

import (
	"slices"

	"github.com/brunosilv96/bs-aesthetics-api/internal/apperrors"
	"github.com/gin-gonic/gin"
)

var rolePermissions = map[string][]string{
	"ADMIN": {
		"CREATE_USER",
		"DELETE_USER",
		"UPDATE_SERVICE",
		"VIEW_REPORTS",
	},
	"CUSTOMER": {
		"CREATE_APPOINTMENT",
		"VIEW_OWN_APPOINTMENTS",
	},
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")

		if userRole == "" {
			c.Error(apperrors.ErrForbidden)
			c.Abort()
			return
		}

		permissions, exists := rolePermissions[userRole]

		if !exists {
			c.Error(apperrors.ErrForbidden)
			c.Abort()
			return
		}

		if slices.Contains(permissions, permission) {
			c.Next()
			return
		}

		c.Error(apperrors.ErrForbidden)
		c.Abort()
	}
}
