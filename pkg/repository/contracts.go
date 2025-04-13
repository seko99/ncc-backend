package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_contracts_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository Contracts

type Contracts interface {
	Create(data models.ContractData) error
	Update(data models.ContractData) error
	Delete(id string) error
	Get() ([]models.ContractData, error)
}
