package database

import (
	"context"

	"github.com/brunosilv96/bs-aesthetics-api/pkg"
	"github.com/jackc/pgx/v5"
)

func Run() error {
	ctx := context.Background()
	cfg := pkg.Load()

	conn, err := pgx.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}

	defer conn.Close(ctx)

	return nil
}
