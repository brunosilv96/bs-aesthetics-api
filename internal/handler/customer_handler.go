package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/brunosilv96/bs-aesthetics-api/internal/exception"
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

func (handler CustomerHandler) RegisterCustomer(c *gin.Context) {
	var payload model.CreateCustomer

	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("error on bind body payload", "error", err)
		exception.BadRequestErrorHandler(c, err)
		return
	}

	customer, err := handler.service.Register(c, payload)
	if err != nil {
		if errors.Is(err, exception.ErrSysCustomerAlreadyRegistered) {
			slog.Error("error on register customer, email already registered", "error", err)
			c.Error(exception.ErrCustomerAlreadyRegistered)
			return
		}

		slog.Error("error on register customer", "error", err)
		c.Error(exception.ErrBadRequest)
		return
	}

	c.JSON(http.StatusCreated, model.CustomerResponse{
		ID:        customer.ID.String(),
		Name:      customer.Name,
		Email:     customer.Email,
		Phone:     customer.Phone,
		Birthdate: customer.Birthdate.Time.Format("2006-01-02"),
		Role:      customer.Role,
		CreatedAt: customer.CreatedAt,
	})
}

func (handler CustomerHandler) Customers(c *gin.Context) {
	customers, err := handler.service.List(c)
	if err != nil {
		slog.Error("Error on load customer list", "error", err)
		c.Error(exception.ErrBadRequest)
		return
	}

	customerList := []model.CustomerResponse{}
	for _, customer := range customers {
		customerList = append(customerList, model.CustomerResponse{
			ID:        customer.ID.String(),
			Name:      customer.Name,
			Email:     customer.Email,
			Phone:     customer.Phone,
			Birthdate: customer.Birthdate.Time.Format("2006-01-02"),
			Role:      customer.Role,
			CreatedAt: customer.CreatedAt,
			UpdatedAt: customer.UpdatedAt,
			DeletedAt: customer.DeletedAt,
		})
	}

	c.JSON(http.StatusOK, customerList)
}

func (handler CustomerHandler) CustomerByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		slog.Error("Parameter ID is empty")
		c.Error(exception.ErrBadRequestID)
		return
	}

	customer, err := handler.service.FindByID(c, id)
	if err != nil {
		slog.Error("Customer not found", "id", id, "error", err)
		if errors.Is(err, exception.ErrSysCustomerNotFound) {
			c.Error(exception.ErrNotFound)
			return
		}

		c.Error(exception.ErrBadRequest)
		return
	}

	c.JSON(http.StatusOK, model.CustomerResponse{
		ID:        customer.ID.String(),
		Name:      customer.Name,
		Email:     customer.Email,
		Phone:     customer.Phone,
		Birthdate: customer.Birthdate.Time.Format("2006-01-02"),
		Role:      customer.Role,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
		DeletedAt: customer.DeletedAt,
	})
}

func (handler CustomerHandler) DisableCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		slog.Error("parameter ID is empty")
		c.Error(exception.ErrBadRequestID)
		return
	}

	err := handler.service.Delete(c, id)
	if err != nil {
		slog.Error("failed to disable customer", "id", id, "error", err)
		c.Error(exception.ErrBadRequest)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (handler CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		slog.Error("parameter ID is empty")
		c.Error(exception.ErrBadRequestID)
		return
	}

	var payload model.UpdateCustomer

	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("error on bind body payload", "error", err)
		exception.BadRequestErrorHandler(c, err)
		return
	}

	if payload.Birthdate != nil {
		_, err := time.Parse("2006-01-02", *payload.Birthdate)
		if err != nil {
			slog.Error("error on bind body payload", "error", err)
			exception.BadRequestErrorHandler(c, err)
			return
		}
	}

	customer, err := handler.service.Update(c, id, payload)
	if err != nil {
		slog.Error("failed to update customer", "id", id, "error", err)
		c.Error(exception.ErrBadRequest)
		return
	}

	c.JSON(http.StatusOK, model.CustomerResponse{
		ID:        customer.ID.String(),
		Name:      customer.Name,
		Email:     customer.Email,
		Phone:     customer.Phone,
		Birthdate: customer.Birthdate.Time.String(),
		Role:      customer.Role,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
		DeletedAt: customer.DeletedAt,
	})
}
