package service

import (
	"context"
	"time"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
	"github.com/brunosilv96/bs-aesthetics-api/internal/model"
	"github.com/brunosilv96/bs-aesthetics-api/internal/repository"
	"github.com/google/uuid"
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

func (service CustomerService) List(ctx context.Context) ([]database.Customer, error) {
	customers, err := service.repository.List(ctx)
	if err != nil {
		return []database.Customer{}, err
	}

	return customers, nil
}

func (service CustomerService) FindByID(ctx context.Context, id string) (database.Customer, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return database.Customer{}, exception.ErrParseUUIDFailed
	}

	customer, err := service.repository.FindByID(ctx, pgtype.UUID{
		Bytes: parsedID,
		Valid: true,
	})
	if err != nil {
		return database.Customer{}, err
	}

	return customer, nil

}

func (service CustomerService) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return exception.ErrParseUUIDFailed
	}

	customer, err := service.repository.FindByID(ctx, pgtype.UUID{
		Bytes: parsedID,
		Valid: true,
	})
	if err != nil {
		return err
	}

	err = service.repository.SoftDelete(ctx, customer.ID)
	if err != nil {
		return err
	}

	return nil

}
