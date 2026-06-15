package router

import (
	"github.com/brunosilv96/bs-aesthetics-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func CustomerRouter(router *gin.Engine, handler handler.CustomerHandler) {
	customer := router.Group("/customer")

	customer.POST("/", handler.Create)
	customer.GET("/", handler.Customers)
	customer.GET("/:id", handler.CustomerByID)
}
