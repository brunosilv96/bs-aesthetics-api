package router

import (
	"github.com/brunosilv96/bs-aesthetics-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func CustomerRouter(router *gin.Engine, handler handler.CustomerHandler) {
	route := router.Group("/customer")

	route.POST("/", handler.RegisterCustomer)
	route.GET("/", handler.Customers)
	route.GET("/:id", handler.CustomerByID)
	route.DELETE("/:id", handler.DisableCustomer)
	route.PATCH("/:id", handler.UpdateCustomer)
}
