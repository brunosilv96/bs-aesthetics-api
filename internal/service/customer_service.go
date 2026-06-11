package service

import (
	"context"
	"time"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
	"github.com/brunosilv96/bs-aesthetics-api/internal/model"
	"github.com/brunosilv96/bs-aesthetics-api/internal/repository"
	"github.com/google/uuid"
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
		ID:        uuid.New(),
		Name:      payload.Name,
		Email:     payload.Email,
		Phone:     payload.Phone,
		BirthDate: payload.BirthDate,
		CreatedAt: time.Now(),
	}

	customer, err := service.repository.Save(ctx, dbModel)
	if err != nil {
		return model.CustomerResponse{}, err
	}

	return model.CustomerResponse{
		ID:        customer.ID,
		Name:      customer.Name,
		Email:     customer.Email,
		Phone:     customer.Phone,
		BirthDate: customer.BirthDate,
		CreatedAt: customer.CreatedAt,
	}, nil
}
