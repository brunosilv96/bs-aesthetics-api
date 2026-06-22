package handler

import (
	"log/slog"
	"net/http"

	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
	"github.com/brunosilv96/bs-aesthetics-api/internal/model"
	"github.com/brunosilv96/bs-aesthetics-api/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (handler *AuthHandler) Login(c *gin.Context) {
	var payload model.AuthRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("Error on bind body payload", "error", err)
		exception.BadRequestErrorHandler(c, err)
		return
	}

	tokens, err := handler.authService.AuthenticateCustomer(c, payload)
	if err != nil {
		slog.Error("Error on generate access for customer", "error", err)
		c.Error(exception.ErrBadRequest)
	}

	c.SetCookie(
		"refresh_token",     // Name
		tokens.RefreshToken, // Value
		7*24*60*60,          // Lifespan (7 days)
		"/",                 // Path (all api)
		"",                  // Domain (empty use the current API domain)
		true,                // Secure (true == HTTPS)
		true,                // HTTP Only (true == prevents javascript from reading)
	)

	// For extra CSRF protection if your front end is web-based:
	// c.Header("SameSite", "Strict")

	c.JSON(http.StatusOK, model.TokenResponse{
		CustomerID:  tokens.CustomerID,
		AccessToken: tokens.AccessToken,
	})
}
