package router

import (
	"github.com/brunosilv96/bs-aesthetics-api/internal/handler"
	"github.com/gin-gonic/gin"
)

func AuthRouter(router *gin.Engine, handler handler.AuthHandler) {
	route := router.Group("/auth")

	route.POST("/login", handler.Login)
	route.POST("/refresh", handler.RefreshToken)
	route.POST("/logoff", handler.Logout)
}
