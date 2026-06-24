package repository

import (
	"context"
	"errors"
	"log/slog"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				slog.Error("[BD] error on save customer, constrants violation", "pgx error:", err)
				return database.Customer{}, exception.ErrSysCustomerAlreadyRegistered
			}
		}

		slog.Error("[BD] error on save customer in database", "pgx error:", err)
		return database.Customer{}, exception.ErrSysSaveCustomer
	}

	return customer, nil
}

func (repository CustomerRepository) List(ctx context.Context) ([]database.Customer, error) {
	customers, err := repository.db.LoadCustomers(ctx)
	if err != nil {
		slog.Error("[BD] error on load customer list in database", "pgx error:", err)
		return []database.Customer{}, exception.ErrSysLoadCustomerList
	}

	return customers, nil
}

func (repository CustomerRepository) FindByID(ctx context.Context, id pgtype.UUID) (database.Customer, error) {
	customer, err := repository.db.FindCustomerByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.Customer{}, exception.ErrSysCustomerNotFound
		}

		slog.Error("[BD] error on load customer by id in database", "id:", id, "pgx error:", err)
		return database.Customer{}, exception.ErrSysFindCustomerById
	}

	return customer, nil
}

func (repository CustomerRepository) FindByEmail(ctx context.Context, email string) (database.Customer, error) {
	customer, err := repository.db.FindCustomerByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.Customer{}, exception.ErrSysCustomerNotFound
		}

		slog.Error("[BD] error on load customer by id in database", "email:", email, "pgx error:", err)
		return database.Customer{}, exception.ErrSysFindCustomerByEmail
	}

	return customer, nil
}

func (repository CustomerRepository) SoftDelete(ctx context.Context, id pgtype.UUID) error {
	err := repository.db.SoftDeleteCustomer(ctx, id)
	if err != nil {
		slog.Error("[BD] error on delete customer by id in database", "id:", id, "pgx error:", err)
		return exception.ErrSysDeleteCustomer
	}

	return nil
}

func (repository CustomerRepository) Update(ctx context.Context, customer database.UpdateCustomerParams) error {
	err := repository.db.UpdateCustomer(ctx, customer)
	if err != nil {
		slog.Error("[BD] error on update customer by id in database", "id:", customer.ID, "pgx error:", err)
		return exception.ErrSysUpdateCustomer
	}

	return nil
}
