package handler

import (
	"errors"
	"log/slog"
	"net/http"

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

	authIdentity, exists := c.Get("auth_identity")
	if !exists {
		slog.Error("error on extract customer_id or role on token")
		c.Error(exception.ErrUnauthorized)
		return
	}

	identity, ok := authIdentity.(*model.Identity)
	if !ok {
		slog.Error("error on extract customer_id or role on token")
		c.Error(exception.ErrUnauthorized)
		return
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("error on bind body payload", "error", err)
		exception.BadRequestErrorHandler(c, err)
		return
	}

	procedure, err := handler.procedureService.Register(c, identity, payload)
	if err != nil {
		if errors.Is(exception.ErrSysNotFound, err) {
			slog.Error("customer id not found", "error: ", err)
			c.Error(exception.ErrNotFound)
			return
		}

		if errors.Is(exception.ErrSysForbidden, err) {
			slog.Error("customer is forbidden or unauthorized", "error: ", err)
			c.Error(exception.ErrForbidden)
			return
		}

		if errors.Is(exception.ErrSysInvalidInput, err) {
			slog.Error("invalid input or use case violate", "error: ", err)
			c.Error(exception.ErrInvalidInput)
			return
		}

		slog.Error("error on register new procedure", "error", err)
		c.Error(exception.ErrInternal)
		return
	}

	c.JSON(http.StatusCreated, model.ProcedureResponse{
		ID:              procedure.ID.String(),
		RegisteredBy:    procedure.RegistredBy.String(),
		Name:            procedure.Name,
		Description:     procedure.Description,
		BannerURL:       &procedure.BannerUrl.String,
		Price:           procedure.Price,
		DurationMinutes: int(procedure.DurationMinutes),
		Available:       procedure.Available,
		CreatedAt:       *procedure.CreatedAt,
	})
}
