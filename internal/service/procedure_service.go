package service

import (
	"context"

	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
	"github.com/brunosilv96/bs-aesthetics-api/internal/model"
	"github.com/brunosilv96/bs-aesthetics-api/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProcedureService struct {
	customerRepository repository.CustomerRepository
}

func NewProcedureService(customerRepository repository.CustomerRepository) *ProcedureService {
	return &ProcedureService{
		customerRepository,
	}
}

func (service *ProcedureService) Register(ctx context.Context, customerID string, payload model.RegisterProcedure) error {
	parsedID, err := uuid.Parse(customerID)
	if err != nil {
		return exception.ErrSysParse
	}

	customer, err := service.customerRepository.FindByID(ctx, pgtype.UUID{
		Bytes: parsedID,
		Valid: true,
	})
	if err != nil {
		return err
	}

	if customer.Role == "customer" {
		return nil
	}

	return nil
}
