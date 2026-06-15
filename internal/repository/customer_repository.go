package repository

import (
	"context"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
)

type CustomerRepository struct {
	db database.Querier
}

func NewCustomerRepository(db database.Querier) *CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (repository CustomerRepository) Save(ctx context.Context, payload database.CreateCustomerParams) (database.Customer, error) {
	customer, err := repository.db.CreateCustomer(ctx, payload)
	if err != nil {
		return database.Customer{}, err
	}

	return customer, nil
}

func (repository CustomerRepository) List(ctx context.Context) ([]database.Customer, error) {
	customers, err := repository.db.LoadCustomers(ctx)
	if err != nil {
		return []database.Customer{}, err
	}

	return customers, nil
}
