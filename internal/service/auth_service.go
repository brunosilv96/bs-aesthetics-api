package service

import (
	"context"
	"time"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
	"github.com/brunosilv96/bs-aesthetics-api/internal/model"
	"github.com/brunosilv96/bs-aesthetics-api/internal/repository"
	"github.com/brunosilv96/bs-aesthetics-api/pkg"
)

type AuthService struct {
	authService     pkg.TokenService
	authRepository  repository.AuthRepository
	customerService CustomerService
}

func NewAuthService(tokenService pkg.TokenService, authRepository repository.AuthRepository, customerService CustomerService) *AuthService {
	return &AuthService{
		authService:     tokenService,
		authRepository:  authRepository,
		customerService: customerService,
	}
}

func (service *AuthService) AuthenticateCustomer(ctx context.Context, payload model.AuthRequest) (*model.TokensOutput, error) {
	customer, err := service.customerService.FindByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}

	isPasswordValid := pkg.CheckPasswordHash(payload.Password, customer.Password)
	if !isPasswordValid {
		return nil, exception.ErrSysInvalidCredential
	}

	tokens, err := service.authService.GenerateToken(customer.ID.String(), customer.Role)
	if err != nil {
		return nil, err
	}

	hashedToken := service.authService.HashRefreshToken(tokens.RefreshToken)

	savedToken, err := service.authRepository.SaveRefreshToken(ctx, database.SaveRefreshTokenParams{
		CustomerID: customer.ID.String(),
		Role:       customer.Role,
		TokenHash:  hashedToken,
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour), // Valid for 1 week
	})
	if err != nil {
		return nil, err
	}

	return &model.TokensOutput{
		CustomerID:   savedToken.CustomerID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (service *AuthService) RotationRefreshToken(ctx context.Context, token string) (*model.TokensOutput, error) {
	hashedToken := service.authService.HashRefreshToken(token)

	foundToken, err := service.authRepository.LoadRefreshToken(ctx, hashedToken)
	if err != nil {
		return nil, err
	}

	if hashedToken != foundToken.TokenHash {
		return nil, exception.ErrSysInvalidToken
	}

	if time.Now().After(foundToken.ExpiresAt) {
		return nil, exception.ErrSysExpiredToken
	}

	newPairToken, err := service.authService.GenerateToken(foundToken.CustomerID, foundToken.Role)
	if err != nil {
		return nil, err
	}

	hashedRefreshToken := service.authService.HashRefreshToken(newPairToken.RefreshToken)

	_, err = service.authRepository.SaveRefreshToken(ctx, database.SaveRefreshTokenParams{
		CustomerID: foundToken.CustomerID,
		Role:       foundToken.Role,
		TokenHash:  hashedRefreshToken,
		ExpiresAt:  time.Now().Add(7 * 24 * time.Hour), // Valid for 1 week
	})
	if err != nil {
		return nil, err
	}

	err = service.authRepository.InvalidRefreshToken(ctx, hashedToken)
	if err != nil {
		return nil, err
	}

	return &model.TokensOutput{
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
		return exception.ErrSysInvalidToken
	}

	err = service.authRepository.InvalidRefreshToken(ctx, hashedToken)
	if err != nil {
		return err
	}

	return nil
}
