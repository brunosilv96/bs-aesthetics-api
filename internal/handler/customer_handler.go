package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/brunosilv96/bs-aesthetics-api/internal/apperrors"
	"github.com/brunosilv96/bs-aesthetics-api/internal/model"
	"github.com/brunosilv96/bs-aesthetics-api/internal/service"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	customerService service.CustomerService
}

func NewCustomerHandler(service service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: service,
	}
}

func (handler CustomerHandler) RegisterCustomer(c *gin.Context) {
	ctx := c.Request.Context()
	var payload model.CreateCustomer

	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("error on bind body payload", "error", err)
		c.Error(
			apperrors.ErrInvalidInput.
				WithDetails(apperrors.PayloadValidator(err)),
		)
		return
	}

	customer, err := handler.customerService.Register(ctx, payload)
	if err != nil {
		if errors.Is(err, apperrors.ErrSysAlreadyExists) {
			slog.Error("error on register customer, email already registered", "error", err)
			c.Error(apperrors.ErrConflict.WithMessage("email already registered"))
			return
		}

		slog.Error("error on register customer", "error", err)
		c.Error(apperrors.ErrUnprocessableEntity)
		return
	}

	c.JSON(http.StatusCreated, &model.CustomerResponse{
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
	ctx := c.Request.Context()

	customers, err := handler.customerService.List(ctx)
	if err != nil {
		slog.Error("Error on load customer list", "error", err)
		c.Error(apperrors.ErrUnprocessableEntity)
		return
	}

	customerList := []model.CustomerResponse{}
	for _, customer := range *customers {
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
	ctx := c.Request.Context()

	id := c.Param("id")
	if id == "" {
		slog.Error("parameter ID is empty")
		c.Error(apperrors.ErrInvalidUUID)
		return
	}

	customer, err := handler.customerService.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			slog.Error("customer not found", "id", id, "error", err)
			c.Error(apperrors.ErrNotFound)
			return
		}

		slog.Error("error on find customer by id", "id", id, "error", err)
		c.Error(apperrors.ErrUnprocessableEntity)
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
		c.Error(apperrors.ErrInvalidUUID)
		return
	}

	err := handler.customerService.Delete(c, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			slog.Error("customer not found", "id", id, "error", err)
			c.Error(apperrors.ErrNotFound)
			return
		}

		slog.Error("failed to disable customer", "id", id, "error", err)
		c.Error(apperrors.ErrUnprocessableEntity)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (handler CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		slog.Error("parameter ID is empty")
		c.Error(apperrors.ErrInvalidUUID)
		return
	}

	var payload model.UpdateCustomer

	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Error("error on bind body payload", "error", err)
		c.Error(
			apperrors.ErrInvalidInput.
				WithDetails(apperrors.PayloadValidator(err)),
		)
		return
	}

	if payload.Birthdate != nil {
		_, err := time.Parse("2006-01-02", *payload.Birthdate)
		if err != nil {
			slog.Error("error on bind body payload", "error", err)
			c.Error(apperrors.ErrInvalidInput)
			return
		}
	}

	customer, err := handler.customerService.Update(c, id, payload)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			slog.Error("customer not found", "id", id, "error", err)
			c.Error(apperrors.ErrNotFound)
			return
		}

		slog.Error("failed to update customer", "id", id, "error", err)
		c.Error(apperrors.ErrBadRequest)
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
