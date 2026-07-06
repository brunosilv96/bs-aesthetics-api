package service

import (
	"context"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
	"github.com/brunosilv96/bs-aesthetics-api/internal/model"
	"github.com/brunosilv96/bs-aesthetics-api/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProcedureService struct {
	customerRepository  repository.CustomerRepository
	procedureRepository repository.ProcedureRepository
}

func NewProcedureService(customerRepository repository.CustomerRepository, procedureRepository repository.ProcedureRepository) *ProcedureService {
	return &ProcedureService{
		customerRepository,
		procedureRepository,
	}
}

func (service *ProcedureService) Register(ctx context.Context, identity *model.Identity, payload model.RegisterProcedure) (*database.Procedure, error) {
	if identity.Role != "admin" {
		return nil, exception.ErrSysForbidden
	}

	if payload.Price < 0 {
		return nil, exception.ErrSysInvalidInput
	}

	parsedID, err := uuid.Parse(identity.CustomerID)
	if err != nil {
		return nil, exception.ErrSysParse
	}

	customer, err := service.customerRepository.FindByID(ctx, pgtype.UUID{
		Bytes: parsedID,
		Valid: true,
	})
	if err != nil {
		return nil, err
	}

	_, err = service.procedureRepository.FindByID(ctx, customer.ID)
	if err != nil && err != exception.ErrSysNotFound {
		return nil, err
	}

	createdProcedure, err := service.procedureRepository.Save(ctx, database.CreateProcedureParams{
		RegistredBy: customer.ID,
		Name:        payload.Name,
		Description: payload.Description,
		BannerUrl: pgtype.Text{
			String: payload.BannerURL,
			Valid:  true,
		},
		Price:           payload.Price,
		DurationMinutes: int32(payload.DurationMinutes),
		Available:       payload.Available,
	})

	return &createdProcedure, nil
}
