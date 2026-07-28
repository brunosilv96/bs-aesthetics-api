package model

import "time"

type EnterpriseRegister struct {
	TradeName   string    `json:"trade_name" binding:"required,min=3"`
	CNPJ        string    `json:"cnpj" binding:"required"`
	OpeningTime time.Time `json:"opening_time" binding:"required"`
	ClosingTime time.Time `json:"closing_time" binding:"required"`
}

type Enterprise struct {
	ID          string     `json:"id"`
	TradeName   string     `json:"trade_name"`
	CNPJ        string     `json:"cnpj"`
	OpeningTime time.Time  `json:"opening_time"`
	ClosingTime time.Time  `json:"closing_time"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}
