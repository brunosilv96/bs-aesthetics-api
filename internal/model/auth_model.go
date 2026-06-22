package model

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	CustomerID  string `json:"customer_id"`
	AccessToken string `json:"access_token"`
}

type TokensOutput struct {
	CustomerID   string
	AccessToken  string
	RefreshToken string
}
