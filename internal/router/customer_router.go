package router

import (
	"github.com/brunosilv96/bs-aesthetics-api/internal/handler"
	"github.com/brunosilv96/bs-aesthetics-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func CustomerRouter(router *gin.Engine, handler handler.CustomerHandler, authMiddleware middleware.AccessTokenMiddleware) {
	route := router.Group("/customer")

	// Public Routes
	route.POST("/", handler.RegisterCustomer)

	route.Use(authMiddleware.RequireAuth())

	// Private Routes
	route.GET("/", handler.Customers)
	route.GET("/:id", handler.CustomerByID)
	route.DELETE("/:id", handler.DisableCustomer)
	route.PATCH("/:id", handler.UpdateCustomer)
}
