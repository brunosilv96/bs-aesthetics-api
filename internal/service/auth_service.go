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
	authRepository  repository.AuthRepository
	customerService CustomerService
}

func NewAuthService(tokenService auth.TokenService, authRepository repository.AuthRepository, customerService CustomerService) *AuthService {
	return &AuthService{
		authService:     tokenService,
		authRepository:  authRepository,
		customerService: customerService,
	}
}

func (service *AuthService) AuthenticateCustomer(ctx context.Context, payload model.AuthRequest) (model.TokensOutput, error) {
	customer, err := service.customerService.FindByEmail(ctx, payload.Email)
	if err != nil {
		return model.TokensOutput{}, exception.ErrSysCustomerNotFound
	}

	isPasswordValid := auth.CheckPasswordHash(payload.Password, customer.Password)
	if !isPasswordValid {
		return model.TokensOutput{}, exception.ErrSysInvalidPassword
	}

	tokens, err := service.authService.GenerateToken(customer.ID.String(), "customer")
	if err != nil {
		return model.TokensOutput{}, err
	}

	hashedToken := service.authService.HashRefreshToken(tokens.RefreshToken)

	savedToken, err := service.authRepository.SaveRefreshToken(ctx, database.SaveRefreshTokenParams{
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

func (service *AuthService) RotationRefreshToken(ctx context.Context, token string) (model.TokensOutput, error) {
	hashedToken := service.authService.HashRefreshToken(token)

	foundToken, err := service.authRepository.LoadRefreshToken(ctx, hashedToken)
	if err != nil {
		return model.TokensOutput{}, err
	}

	if hashedToken != foundToken.TokenHash {
		return model.TokensOutput{}, exception.ErrSysRefreshTokenNotFound
	}

	if time.Now().After(foundToken.ExpiresAt) {
		return model.TokensOutput{}, exception.ErrSysRefreshTokenExpired
	}

	newPairToken, err := service.authService.GenerateToken(foundToken.CustomerID, "customer")
	if err != nil {
		return model.TokensOutput{}, err
	}

	hashedRefreshToken := service.authService.HashRefreshToken(newPairToken.RefreshToken)

	_, err = service.authRepository.SaveRefreshToken(ctx, database.SaveRefreshTokenParams{
		CustomerID: foundToken.ID.String(),
		TokenHash:  hashedRefreshToken,
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour), // Valid for 1 week
	})
	if err != nil {
		return model.TokensOutput{}, err
	}

	err = service.authRepository.InvalidRefreshToken(ctx, hashedToken)
	if err != nil {
		return model.TokensOutput{}, err
	}

	return model.TokensOutput{
		CustomerID:   foundToken.CustomerID,
		AccessToken:  newPairToken.AccessToken,
		RefreshToken: newPairToken.RefreshToken,
	}, nil
}

func (service *AuthService) Logoff(ctx context.Context, token string) error {
	hashedToken := service.authService.HashRefreshToken(token)

	foundToken, err := service.authRepository.LoadRefreshToken(ctx, hashedToken)
	if err != nil {
		return err
	}

	if hashedToken != foundToken.TokenHash {
		return exception.ErrSysRefreshTokenNotFound
	}

	err = service.authRepository.InvalidRefreshToken(ctx, hashedToken)
	if err != nil {
		return err
	}

	return nil
}
