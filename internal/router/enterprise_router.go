package router

import (
	"github.com/brunosilv96/bs-aesthetics-api/internal/handler"
	"github.com/brunosilv96/bs-aesthetics-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func EnterpriseRouter(
	router *gin.Engine,
	handler handler.EnterpriseHandler,
	authMiddleware middleware.AuthMiddleware,
) {
	route := router.Group("/enterprise")

	// Private Routes
	route.Use(authMiddleware.RequireAuth())
	route.Use(authMiddleware.RequireRole("admin"))
	route.POST("/", handler.Register)
}
