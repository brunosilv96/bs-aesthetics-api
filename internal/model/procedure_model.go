package model

import "time"

type RegisterProcedure struct {
	Name            string  `json:"name" binding:"required,min=3"`
	Description     string  `json:"description" binding:"required,min=3"`
	BannerURL       string  `json:"banner_url"`
	Price           float64 `json:"price" binding:"required"`
	DurationMinutes int     `json:"duration_minutes" binding:"required"`
	Available       bool    `json:"available" binding:"required"`
}

type ProcedureResponse struct {
	ID              string     `json:"id"`
	RegisteredBy    string     `json:"registered_by"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	BannerURL       *string    `json:"banner_url,omitempty"`
	Price           float64    `json:"price"`
	DurationMinutes int        `json:"duration_minutes"`
	Available       bool       `json:"available"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}
