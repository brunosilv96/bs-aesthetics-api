package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenService struct {
	jwtSecret      []byte
	issuer         string
	accessDuration time.Duration
}

func NewTokenService(secret, issuer string) *TokenService {
	return &TokenService{
		jwtSecret:      []byte(secret),
		issuer:         issuer,
		accessDuration: 15 * time.Minute, // Time default for access token
	}
}

func (service *TokenService) GenerateToken(customerID, role string) (*TokenPair, error) {
	// 1. Generate Access Token (JWT)
	claims := CustomClaim{
		CustomerID: customerID,
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(service.accessDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    service.issuer,
		},
	}

	// Access Token
	jwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Access Token Assigned
	signedAccessToken, err := jwtAccessToken.SignedString(service.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("error to assign access token: %w", err)
	}

	// 2. Generate Refresh Token
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("error to generate random string: %w", err)
	}

	refreshTokenString := hex.EncodeToString(b)

	return &TokenPair{
		AccessToken:  signedAccessToken,
		RefreshToken: refreshTokenString,
	}, nil
}

func (service *TokenService) HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))

	return hex.EncodeToString(hash[:])
}
