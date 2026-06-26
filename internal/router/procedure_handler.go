package router

import (
	"github.com/brunosilv96/bs-aesthetics-api/internal/handler"
	"github.com/brunosilv96/bs-aesthetics-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ProcedureRouter(router *gin.Engine, handler handler.ProcedureHandler, authMiddleware middleware.AccessTokenMiddleware) {
	route := router.Group("/procedure")

	route.POST("/register", handler.Create)
}
