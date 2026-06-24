package service

import (
	"context"
	"time"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
	"github.com/brunosilv96/bs-aesthetics-api/internal/auth"
	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
	"github.com/brunosilv96/bs-aesthetics-api/internal/model"
	"github.com/brunosilv96/bs-aesthetics-api/internal/repository"
)

type AuthService struct {
	authService     auth.TokenService
	AuthRepository  repository.AuthRepository
	customerService CustomerService
}

func NewAuthService(tokenService auth.TokenService, authRepository repository.AuthRepository, customerService CustomerService) *AuthService {
	return &AuthService{
		authService:     tokenService,
		AuthRepository:  authRepository,
		customerService: customerService,
	}
}

func (service *AuthService) AuthenticateCustomer(ctx context.Context, payload model.AuthRequest) (model.TokensOutput, error) {
	customer, err := service.customerService.FindByEmail(ctx, payload.Email)
	if err != nil {
		return model.TokensOutput{}, exception.ErrSysCustomerNotFound
	}

	passwordNonMatches := auth.CheckPasswordHash(payload.Password, customer.Password)
	if passwordNonMatches {
		return model.TokensOutput{}, exception.ErrSysInvalidPassword
	}

	tokens, err := service.authService.GenerateToken(customer.ID.String(), "customer")
	if err != nil {
		return model.TokensOutput{}, err
	}

	hashedToken := service.authService.HashRefreshToken(tokens.RefreshToken)

	savedToken, err := service.AuthRepository.SaveRefreshToken(ctx, database.SaveRefreshTokenParams{
		CustomerID: customer.ID.String(),
		TokenHash:  hashedToken,
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour), // Valid for 1 week
	})
	if err != nil {
		return model.TokensOutput{}, err
	}

	return model.TokensOutput{
		CustomerID:   savedToken.CustomerID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
