package pkg

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
	"github.com/brunosilv96/bs-aesthetics-api/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenService struct {
	jwtSecret           []byte
	issuer              string
	accessTokenDuration time.Duration
}

func NewTokenService(secret, issuer string) *TokenService {
	return &TokenService{
		jwtSecret:           []byte(secret),
		issuer:              issuer,
		accessTokenDuration: 15 * time.Minute, // Time default for access token
	}
}

func (service *TokenService) GenerateToken(customerID, role string) (*TokenPair, error) {
	// 1. Generate Access Token (JWT)
	claims := model.CustomClaim{
		CustomerID: customerID,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(service.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    service.issuer,
		},
	}

	// Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Access Token Assigned
	signedAccessToken, err := accessToken.SignedString(service.jwtSecret)
	if err != nil {
		slog.Error("error on assigning access token", "error:", err)
		return nil, fmt.Errorf("%v - Error: %w", exception.ErrSysGenerate, err)
	}

	// 2. Generate Refresh Token
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		slog.Error("failed on generate random bytes for refresh token", "error:", err)
		return nil, fmt.Errorf("%v - Error: %w", exception.ErrSysGenerate, err)
	}

	refreshToken := hex.EncodeToString(b)

	return &TokenPair{
		AccessToken:  signedAccessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (service *TokenService) HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))

	return hex.EncodeToString(hash[:])
}

func (service *TokenService) ValidateAccessToken(encryptedToken string) (*model.CustomClaim, error) {
	token, err := jwt.Parse(encryptedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Error("error on verify token assigning, most different")
			return nil, exception.ErrSysInvalidToken
		}

		return []byte(service.jwtSecret), nil
	})
	if err != nil {
		slog.Error("failed to parse bearer token", "error", err)
		return nil, exception.ErrSysInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		slog.Error("bearer token is invalid or not match with custom claim", "error", err)
		return nil, exception.ErrSysInvalidToken
	}

	return &model.CustomClaim{
		CustomerID: claims["customer_id"].(string),
		Role:       claims["role"].(string),
	}, nil
}
