package main

import (
	"github.com/brunosilv96/bs-aesthetics-api/internal/router"
	"github.com/brunosilv96/bs-aesthetics-api/pkg"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configurations
	cfg := pkg.Load()

	gin.SetMode(gin.DebugMode)

	if cfg.Environment == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	if cfg.Environment == "test" {
		gin.SetMode(gin.TestMode)
	}

	r := gin.New()

	// Middleware
	r.Use(gin.Logger())

	router.NewRouter(r)

	r.Run()
}
