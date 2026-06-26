package repository

import (
	"context"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
)

type ProcedureRepository struct {
	db database.Querier
}

func NewProcedureRepository(db database.Querier) *ProcedureRepository {
	return &ProcedureRepository{
		db,
	}
}

func (repository *ProcedureRepository) Save(ctx context.Context) error {
	return nil
}
