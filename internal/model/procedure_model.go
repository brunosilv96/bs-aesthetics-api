package model

type RegisterProcedure struct {
	Name            string  `json:"name" binding:"required,min=3"`
	Description     string  `json:"description" binding:"required,min=3"`
	BannerURL       string  `json:"banner_url"`
	Price           float64 `json:"price" binding:"required"`
	DurationMinutes int     `json:"duration_minutes" binding:"required"`
	Available       bool    `json:"available" binding:"required"`
}
