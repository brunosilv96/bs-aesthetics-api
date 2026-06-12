package service

import (
	"context"
	"time"

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

func (service CustomerService) Register(ctx context.Context, payload model.CreateCustomer) (database.Customer, error) {
	birthdate, err := time.Parse("2006-01-02", payload.Birthdate)
	if err != nil {
		return database.Customer{}, err
	}

	dbModel := database.CreateCustomerParams{
		Name:     payload.Name,
		Email:    payload.Email,
		Phone:    payload.Phone,
		Password: payload.Password,
		Birthdate: pgtype.Date{
			Time:  birthdate,
			Valid: true,
		},
	}

	customer, err := service.repository.Save(ctx, dbModel)
	if err != nil {
		return database.Customer{}, err
	}

	return customer, nil
}
