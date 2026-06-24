package repository

import (
	"context"
	"errors"
	"log/slog"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
	"github.com/jackc/pgx/v5"
)

type AuthRepository struct {
	db database.Querier
}

func NewAuthRepository(db database.Querier) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (repository *AuthRepository) SaveRefreshToken(ctx context.Context, payload database.SaveRefreshTokenParams) (database.RefreshToken, error) {
	refreshToken, err := repository.db.SaveRefreshToken(ctx, payload)
	if err != nil {
		slog.Error("[BD] error on save refresh token in database", "pgx error:", err)
		return database.RefreshToken{}, exception.ErrSysSaveRefreshToken
	}

	return refreshToken, nil
}

func (repository *AuthRepository) LoadRefreshToken(ctx context.Context, hashToken string) (database.RefreshToken, error) {
	refreshToken, err := repository.db.FindRefreshToken(ctx, hashToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.RefreshToken{}, exception.ErrSysRefreshTokenNotFound
		}

		slog.Error("[BD] error on find refresh token in database", "pgx error:", err)
		return database.RefreshToken{}, exception.ErrSysFindRefreshToken
	}

	return refreshToken, nil
}

func (repository *AuthRepository) InvalidRefreshToken(ctx context.Context, hashToken string) error {
	err := repository.db.InvalidRefreshToken(ctx, hashToken)
	if err != nil {
		slog.Error("[BD] error on find refresh token in database", "pgx error:", err)
		return exception.ErrSysInvalidRefreshToken
	}

	return nil
}
