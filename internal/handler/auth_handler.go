package handler

import (
	"errors"
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
		slog.Error("error on bind body payload", "error", err)
		exception.BadRequestErrorHandler(c, err)
		return
	}

	tokens, err := handler.authService.AuthenticateCustomer(c, payload)
	if err != nil {
		if errors.Is(err, exception.ErrSysCustomerNotFound) {
			slog.Error("customer not found", "email", payload.Email, "error", err)
			c.Error(exception.ErrCustomerNotFound)
			return
		}

		if errors.Is(err, exception.ErrSysInvalidPassword) {
			slog.Error("customer password is invalid", "email", payload.Email, "error", err)
			c.Error(exception.ErrInvalidPassword)
			return
		}

		slog.Error("error on generate access for customer", "error", err)
		c.Error(exception.ErrBadRequest)
		return
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
	c.Header("SameSite", "Strict")

	c.JSON(http.StatusOK, model.TokenResponse{
		CustomerID:  tokens.CustomerID,
		AccessToken: tokens.AccessToken,
	})
}

func (handler *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		slog.Error("error recovery refresh token cookie", "error", err)
		c.Error(exception.ErrRefreshTokenBadRequest)
		return
	}

	tokens, err := handler.authService.RotationRefreshToken(c, refreshToken)
	if err != nil {
		if errors.Is(err, exception.ErrSysRefreshTokenNotFound) {
			slog.Error("refresh token not found", "error", err)
			c.Error(exception.ErrRefreshTokenBadRequest)
			return
		}

		if errors.Is(err, exception.ErrSysRefreshTokenExpired) {
			slog.Error("refresh token has expired", "error", err)
			c.Error(exception.ErrRefreshTokenExpired)
			return
		}

		slog.Error("error rotate refresh token", "error", err)
		c.Error(exception.ErrRefreshTokenBadRequest)
		return
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
	c.Header("SameSite", "Strict")

	c.JSON(http.StatusOK, model.TokenResponse{
		CustomerID:  tokens.CustomerID,
		AccessToken: tokens.AccessToken,
	})
}

func (handler *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{})
		return
	}

	err = handler.authService.Logoff(c, refreshToken)
	if err != nil {
		if errors.Is(err, exception.ErrSysRefreshTokenNotFound) {
			slog.Error("refresh token not found", "error", err)
			c.Error(exception.ErrRefreshTokenBadRequest)
			return
		}

		if errors.Is(err, exception.ErrSysRefreshTokenExpired) {
			slog.Error("refresh token has expired", "error", err)
			c.Error(exception.ErrRefreshTokenExpired)
			return
		}

		slog.Error("error on invalidate refresh token", "error", err)
		c.Error(exception.ErrInvalidateRefreshToken)
		return
	}

	c.SetCookie(
		"refresh_token", // Name
		"",              // Value
		-1,              // Lifespan
		"/",             // Path (all api)
		"",              // Domain (empty use the current API domain)
		true,            // Secure (true == HTTPS)
		true,            // HTTP Only (true == prevents javascript from reading)
	)

	// For extra CSRF protection if your front end is web-based:
	c.Header("SameSite", "Strict")

	c.JSON(http.StatusNoContent, gin.H{})
}
