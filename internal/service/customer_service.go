package service

import (
	"context"
	"time"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
	"github.com/brunosilv96/bs-aesthetics-api/internal/apperrors"
	"github.com/brunosilv96/bs-aesthetics-api/internal/model"
	"github.com/brunosilv96/bs-aesthetics-api/internal/repository"
	"github.com/brunosilv96/bs-aesthetics-api/pkg"
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

func (service CustomerService) Register(ctx context.Context, payload model.CreateCustomer) (*database.Customer, error) {
	foundCustomer, err := service.repository.FindByEmail(ctx, payload.Email)
	if err != nil && err != apperrors.ErrSysNotFound {
		return nil, err
	}

	if foundCustomer != nil {
		return nil, apperrors.ErrSysAlreadyExists
	}

	hashedPassword, err := pkg.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	birthdate, err := time.Parse("2006-01-02", payload.Birthdate)
	if err != nil {
		return nil, apperrors.ErrSysParse
	}

	customer, err := service.repository.Save(ctx, database.CreateCustomerParams{
		Name:     payload.Name,
		Email:    payload.Email,
		Phone:    payload.Phone,
		Password: hashedPassword,
		Role:     payload.Role,
		Birthdate: pgtype.Date{
			Time:  birthdate,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (service CustomerService) List(ctx context.Context) (*[]database.Customer, error) {
	customers, err := service.repository.List(ctx)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (service CustomerService) FindByID(ctx context.Context, id string) (*database.Customer, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperrors.ErrSysParse
	}

	customer, err := service.repository.FindByID(ctx, pgtype.UUID{
		Bytes: parsedID,
		Valid: true,
	})
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (service CustomerService) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return apperrors.ErrSysParse
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

func (service CustomerService) Update(ctx context.Context, id string, payload model.UpdateCustomer) (*database.Customer, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperrors.ErrSysParse
	}

	customer, err := service.repository.FindByID(ctx, pgtype.UUID{
		Bytes: parsedID,
		Valid: true,
	})
	if err != nil {
		return nil, err
	}

	if payload.Name != nil {
		customer.Name = *payload.Name
	}

	if payload.Email != nil {
		customer.Email = *payload.Email
	}

	if payload.Phone != nil {
		customer.Phone = *payload.Phone
	}

	if payload.Role != nil {
		customer.Role = *payload.Role
	}

	if payload.Birthdate != nil {
		birthdate, err := time.Parse("2006-01-02", *payload.Birthdate)
		if err != nil {
			return nil, apperrors.ErrSysParse
		}

		customer.Birthdate = pgtype.Date{
			Time:  birthdate,
			Valid: true,
		}
	}

	err = service.repository.Update(ctx, database.UpdateCustomerParams{
		ID:        customer.ID,
		Name:      customer.Name,
		Email:     customer.Email,
		Phone:     customer.Phone,
		Birthdate: customer.Birthdate,
		Role:      customer.Role,
	})
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (service CustomerService) FindByEmail(ctx context.Context, email string) (*database.Customer, error) {
	customer, err := service.repository.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return customer, nil
}
