package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_hardware_models_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository HardwareModels

type HardwareModels interface {
	Create(data models.HardwareModelData) error
	Upsert(data models.HardwareModelData) error
	Update(data models.HardwareModelData) error
	Delete(id string) error
	Get() ([]models.HardwareModelData, error)
}
