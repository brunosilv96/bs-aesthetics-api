package repository

import (
	"context"
	"log/slog"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
)

type EnterpriseRepository struct {
	db database.Querier
}

func NewEnterpriseRepository(db database.Querier) *EnterpriseRepository {
	return &EnterpriseRepository{
		db,
	}
}

func (repository *EnterpriseRepository) Save(ctx context.Context, payload *database.RegisterEnterpriseParams) (*database.Enterprise, error) {
	enterprise, err := repository.db.RegisterEnterprise(ctx, *payload)
	if err != nil {
		slog.Error("[BD] error on load procedure by id in database", "pgx error:", err)
		return nil, err
	}

	return &enterprise, nil
}
