package service

import (
	"context"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
	"github.com/brunosilv96/bs-aesthetics-api/internal/model"
	"github.com/brunosilv96/bs-aesthetics-api/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type CustomerService struct {
	repository repository.CustomerRepository
}

func NewCustomerService(repository repository.CustomerRepository) *CustomerService {
	return &CustomerService{
		repository: repository,
	}
}

func (service CustomerService) Register(ctx context.Context, payload model.CreateCustomer) (model.CustomerResponse, error) {
	dbModel := database.CreateCustomerParams{
		Name:     payload.Name,
		Email:    payload.Email,
		Phone:    payload.Phone,
		Password: payload.Password,
		BirthDate: pgtype.Timestamp{
			Time:  payload.BirthDate,
			Valid: true,
		},
	}

	customer, err := service.repository.Save(ctx, dbModel)
	if err != nil {
		return model.CustomerResponse{}, err
	}

	return model.CustomerResponse{
		ID:        customer.ID.String(),
		Name:      customer.Name,
		Email:     customer.Email,
		Phone:     customer.Phone,
		BirthDate: customer.BirthDate.Time,
		CreatedAt: customer.CreatedAt.Time,
	}, nil
}
