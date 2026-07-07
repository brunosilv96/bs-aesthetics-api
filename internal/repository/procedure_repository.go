package repository

import (
	"context"
	"errors"
	"log/slog"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
	"github.com/brunosilv96/bs-aesthetics-api/internal/apperrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProcedureRepository struct {
	db database.Querier
}

func NewProcedureRepository(db database.Querier) *ProcedureRepository {
	return &ProcedureRepository{
		db,
	}
}

func (repository *ProcedureRepository) Save(ctx context.Context, payload database.CreateProcedureParams) (database.Procedure, error) {
	procedure, err := repository.db.CreateProcedure(ctx, payload)
	if err != nil {
		slog.Error("[BD] error on load procedure by id in database", "pgx error:", err)
		return database.Procedure{}, err
	}

	return procedure, err
}

func (repository *ProcedureRepository) FindByID(ctx context.Context, id pgtype.UUID) (database.Procedure, error) {
	procedure, err := repository.db.FindProcedureByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.Procedure{}, apperrors.ErrSysNotFound
		}

		slog.Error("[BD] error on load procedure by id in database", "id:", id, "pgx error:", err)
		return database.Procedure{}, err
	}

	return procedure, nil
}
