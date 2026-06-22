package auth

import "github.com/golang-jwt/jwt/v5"

type CustomClaim struct {
	CustomerID string `json:"customer_id"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}
