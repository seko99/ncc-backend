package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_sorm_gateway_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository SormGateway

type SormGateway interface {
	Create(data models.SormGatewayData) error
	Upsert(data models.SormGatewayData) error
	Update(data models.SormGatewayData) error
	Delete(id string) error
	Get() ([]models.SormGatewayData, error)
}
