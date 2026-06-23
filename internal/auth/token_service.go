package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
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
	claims := CustomClaim{
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
		return nil, fmt.Errorf("error to assign access token: %w", err)
	}

	// 2. Generate Refresh Token
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("error to generate random string: %w", err)
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

func (service *TokenService) ValidateAccessToken(encryptedToken string) (*CustomClaim, error) {
	token, err := jwt.Parse(encryptedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, exception.ErrSysTokenAssigningInvalid
		}

		return []byte(service.jwtSecret), nil
	})
	if err != nil {
		slog.Error("failed to parse bearer token", "error", err)
		return nil, exception.ErrSysExpiredOrInvalidBearerToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		slog.Error("bearer token is invalid or not match with custom claim", "error", err)
		return nil, exception.ErrSysExpiredOrInvalidBearerToken
	}

	return &CustomClaim{
		CustomerID: claims["customer_id"].(string),
		Role:       claims["role"].(string),
	}, nil
}
