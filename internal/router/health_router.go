package router

import (
	"github.com/brunosilv96/bs-aesthetics-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func HealthCheckRouter(router *gin.Engine) {
	router.GET("/health", handler.HealthCheck)
}
