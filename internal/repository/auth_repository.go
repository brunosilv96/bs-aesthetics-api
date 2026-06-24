package repository

import (
	"context"
	"log/slog"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
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
	refresh_token, err := repository.db.SaveRefreshToken(ctx, payload)
	if err != nil {
		slog.Error("[BD] error on save refresh token in database", "pgx error:", err)
		return database.RefreshToken{}, exception.ErrSysSaveRefreshToken
	}

	return refresh_token, nil
}
