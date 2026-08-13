package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/brunosilv96/bs-aesthetics-api/config"
	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
	"github.com/brunosilv96/bs-aesthetics-api/internal/handler"
	"github.com/brunosilv96/bs-aesthetics-api/internal/middleware"
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
	tokenService := pkg.NewTokenService(cfg.JwtSecret, "bs-aesthetics-api")

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
	procedureRepository := repository.NewProcedureRepository(store)
	enterpriseRepository := repository.NewEnterpriseRepository(store)

	customerService := service.NewCustomerService(*customerRepository)
	authService := service.NewAuthService(*tokenService, *authRepository, *customerService)
	procedureService := service.NewProcedureService(*customerRepository, *procedureRepository)
	enterpriseService := service.NewEnterpriseService(*enterpriseRepository)

	customerHandler := handler.NewCustomerHandler(*customerService)
	authHandler := handler.NewAuthHandler(*authService)
	procedureHandler := handler.NewProcedureHandler(*procedureService)
	enterpriseHandler := handler.NewEnterpriseHandler(*enterpriseService)

	// Server Config
	r := gin.New()

	r.SetTrustedProxies(nil)

	// Middlewares
	r.Use(gin.Logger())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.TimeoutMiddleware(10 * time.Second))
	authMiddleware := middleware.NewAccessTokenMiddleware(*tokenService)

	router.HealthCheckRouter(r)
	router.CustomerRouter(r, *customerHandler, *authMiddleware)
	router.AuthRouter(r, *authHandler)
	router.ProcedureRouter(r, *procedureHandler, *authMiddleware)
	router.EnterpriseRouter(r, *enterpriseHandler, *authMiddleware)

	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	slog.Info("HTTP server started", "addr", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Error to initialize server: ", err)
	}
}
