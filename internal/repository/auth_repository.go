package repository

import (
	"context"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
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
		return database.RefreshToken{}, err
	}

	return refresh_token, nil
}
