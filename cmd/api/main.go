package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/brunosilv96/bs-aesthetics-api/config"
	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
	"github.com/brunosilv96/bs-aesthetics-api/internal/auth"
	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
	"github.com/brunosilv96/bs-aesthetics-api/internal/handler"
	"github.com/brunosilv96/bs-aesthetics-api/internal/repository"
	"github.com/brunosilv96/bs-aesthetics-api/internal/router"
	"github.com/brunosilv96/bs-aesthetics-api/internal/service"
	"github.com/brunosilv96/bs-aesthetics-api/pkg"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	// Load configurations
	cfg := config.Load()
	logger := pkg.New()

	slog.SetDefault(logger)
	tokenService := auth.NewTokenService(cfg.JwtSecret, "bs-aesthetics-api")

	slog.Info("Server starting...")

	// Start Database Connection
	conn, err := pgx.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Error to initialize database: ", err)
	}
	defer conn.Close(ctx)

	store := database.New(conn)

	// Dependency Injection
	customerRepository := repository.NewCustomerRepository(store)
	authRepository := repository.NewAuthRepository(store)

	customerService := service.NewCustomerService(*customerRepository)
	authService := service.NewAuthService(*tokenService, *authRepository, *customerService)

	customerHandler := handler.NewCustomerHandler(*customerService)
	authHandler := handler.NewAuthHandler(*authService)

	// Server Config
	r := gin.New()

	r.SetTrustedProxies(nil)

	// Middlewares
	r.Use(gin.Logger())
	r.Use(exception.ErrorHandler())

	router.HealthCheckRouter(r)
	router.CustomerRouter(r, *customerHandler)
	router.AuthRouter(r, *authHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Error to initialize server: ", err)
	}
}
