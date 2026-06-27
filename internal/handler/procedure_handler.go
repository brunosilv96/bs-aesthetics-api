package handler

import (
	"log/slog"

	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
	"github.com/brunosilv96/bs-aesthetics-api/internal/model"
	"github.com/brunosilv96/bs-aesthetics-api/internal/service"
	"github.com/gin-gonic/gin"
)

type ProcedureHandler struct {
	procedureService service.ProcedureService
}

func NewProcedureHandler(procedureService service.ProcedureService) *ProcedureHandler {
	return &ProcedureHandler{
		procedureService: procedureService,
	}
}

func (handler *ProcedureHandler) Create(c *gin.Context) {
	var payload model.RegisterProcedure
	customerID, err := c.Cookie("customer_id")
	if err != nil {
		slog.Error("error on extract customer_id on cookie", "error", err)
		c.Error(exception.ErrUnauthorized)
		return
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("error on bind body payload", "error", err)
		exception.BadRequestErrorHandler(c, err)
		return
	}

	err = handler.procedureService.Register(c, customerID, payload)
	if err != nil {
		slog.Error("error on register new procedure", "error", err)
		c.Error(exception.ErrInternal)
		return
	}

	// c.JSON(http.StatusCreated, procedure)
}
