package handler

import (
	"log/slog"
	"net/http"

	"github.com/brunosilv96/bs-aesthetics-api/internal/error"
	"github.com/brunosilv96/bs-aesthetics-api/internal/model"
	"github.com/brunosilv96/bs-aesthetics-api/internal/service"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service service.CustomerService
}

func NewCustomerHandler(service service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		service: service,
	}
}

func (handler CustomerHandler) Create(c *gin.Context) {
	var payload model.CreateCustomer

	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("Error on bind body payload", "error", err)
		error.BadRequestErrorHandler(c, err)
		return
	}

	customer, err := handler.service.Register(c, payload)
	if err != nil {
		slog.Error("Error on register customer", "error", err)
		c.Error(error.ErrBadRequest)
		return
	}

	c.JSON(http.StatusCreated, model.CustomerResponse{
		ID:        customer.ID.String(),
		Name:      customer.Name,
		Email:     customer.Email,
		Phone:     customer.Phone,
		Birthdate: customer.Birthdate.Time.Format("2006-01-02"),
		CreatedAt: customer.CreatedAt.Time,
	})
}
