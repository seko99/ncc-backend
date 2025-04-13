package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_devices_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository Devices

type Devices interface {
	Create(data models.DeviceData) error
	Update(data models.DeviceData) error
	Delete(id string) error
	Get() ([]models.DeviceData, error)
}
