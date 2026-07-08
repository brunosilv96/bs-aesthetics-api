package model

import (
	"time"
)

type CreateCustomer struct {
	Name      string `json:"name" binding:"required,min=3"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Birthdate string `json:"birthdate" binding:"required,datetime=2006-01-02"`
}
type UpdateCustomer struct {
	Name      *string `json:"name,omitempty" validate:"min=3"`
	Email     *string `json:"email,omitempty" validate:"email"`
	Phone     *string `json:"phone,omitempty"`
	Birthdate *string `json:"birthdate,omitempty" validate:"datetime=2006-01-02"`
}

type CustomerResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Birthdate string     `json:"birthdate"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
