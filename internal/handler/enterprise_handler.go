package handler

import (
	"log/slog"
	"net/http"

	"github.com/brunosilv96/bs-aesthetics-api/internal/apperrors"
	"github.com/brunosilv96/bs-aesthetics-api/internal/model"
	"github.com/brunosilv96/bs-aesthetics-api/internal/service"
	"github.com/gin-gonic/gin"
)

type EnterpriseHandler struct {
	enterpriseService service.EnterpriseService
}

func NewEnterpriseHandler(enterpriseService service.EnterpriseService) *EnterpriseHandler {
	return &EnterpriseHandler{
		enterpriseService,
	}
}

func (handler *EnterpriseHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var payload model.EnterpriseRegister

	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("error on bind body payload", "error", err)
		c.Error(
			apperrors.ErrInvalidInput.
				WithDetails(apperrors.PayloadValidator(err)),
		)
		return
	}

	enterprise, err := handler.enterpriseService.Register(ctx, payload)
	if err != nil {
		slog.Error("error on register the enterprise", "error", err)
		c.Error(apperrors.ErrUnprocessableEntity)
		return
	}

	c.JSON(http.StatusCreated, enterprise)
}
