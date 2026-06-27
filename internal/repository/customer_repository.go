package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
		slog.Error("[BD] error on save customer", "error:", err)
		return database.Customer{}, fmt.Errorf("%v - DB Error: %w", exception.ErrSysInternal, err)
	}

	return customer, nil
}

func (repository CustomerRepository) List(ctx context.Context) ([]database.Customer, error) {
	customers, err := repository.db.LoadCustomers(ctx)
	if err != nil {
		slog.Error("[BD] error on load customer list in database", "pgx error:", err)
		return []database.Customer{}, fmt.Errorf("%v - DB Error: %w", exception.ErrSysInternal, err)
	}

	return customers, nil
}

func (repository CustomerRepository) FindByID(ctx context.Context, id pgtype.UUID) (database.Customer, error) {
	customer, err := repository.db.FindCustomerByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.Customer{}, exception.ErrSysNotFound
		}

		slog.Error("[BD] error on load customer by id in database", "id:", id, "pgx error:", err)
		return database.Customer{}, fmt.Errorf("%v - DB Error: %w", exception.ErrSysInternal, err)
	}

	return customer, nil
}

func (repository CustomerRepository) FindByEmail(ctx context.Context, email string) (database.Customer, error) {
	customer, err := repository.db.FindCustomerByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.Customer{}, exception.ErrSysNotFound
		}

		slog.Error("[BD] error on load customer by id in database", "email:", email, "pgx error:", err)
		return database.Customer{}, fmt.Errorf("%v - DB Error: %w", exception.ErrSysInternal, err)
	}

	return customer, nil
}

func (repository CustomerRepository) SoftDelete(ctx context.Context, id pgtype.UUID) error {
	err := repository.db.SoftDeleteCustomer(ctx, id)
	if err != nil {
		slog.Error("[BD] error on delete customer by id in database", "id:", id, "pgx error:", err)
		return fmt.Errorf("%v - DB Error: %w", exception.ErrSysInternal, err)
	}

	return nil
}

func (repository CustomerRepository) Update(ctx context.Context, customer database.UpdateCustomerParams) error {
	err := repository.db.UpdateCustomer(ctx, customer)
	if err != nil {
		slog.Error("[BD] error on update customer by id in database", "id:", customer.ID, "pgx error:", err)
		return fmt.Errorf("%v - DB Error: %w", exception.ErrSysInternal, err)
	}

	return nil
}
