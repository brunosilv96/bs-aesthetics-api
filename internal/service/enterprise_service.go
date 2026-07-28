package service

import (
	"context"

	database "github.com/brunosilv96/bs-aesthetics-api/database/sqlc"
	"github.com/brunosilv96/bs-aesthetics-api/internal/model"
	"github.com/brunosilv96/bs-aesthetics-api/internal/repository"
)

type EnterpriseService struct {
	enterpriseRepository repository.EnterpriseRepository
}

func NewEnterpriseService(enterpriseRepository repository.EnterpriseRepository) *EnterpriseService {
	return &EnterpriseService{
		enterpriseRepository,
	}
}

func (service *EnterpriseService) Register(
	ctx context.Context,
	payload model.EnterpriseRegister,
) (*model.Enterprise, error) {
	enterprise, err := service.enterpriseRepository.Save(ctx, &database.RegisterEnterpriseParams{
		TradeName:   payload.TradeName,
		Cnpj:        payload.CNPJ,
		OpeningTime: payload.OpeningTime,
		ClosingTime: payload.ClosingTime,
	})
	if err != nil {
		return nil, err
	}

	return &model.Enterprise{
		ID:          enterprise.ID.String(),
		TradeName:   enterprise.TradeName,
		CNPJ:        enterprise.Cnpj,
		OpeningTime: enterprise.OpeningTime.Local().UTC(),
		ClosingTime: enterprise.ClosingTime.Local().UTC(),
		CreatedAt:   enterprise.CreatedAt.Local().UTC(),
		UpdatedAt:   enterprise.UpdatedAt,
	}, nil
}
