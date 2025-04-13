package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_device_interfaces_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository DeviceInterfaces

type DeviceInterfaces interface {
	Create(data models.IfaceData) error
	Update(data models.IfaceData) error
	Delete(id string) error
	Get() ([]models.IfaceData, error)
	GetByDeviceIdAndPort(id string, port int) (models.IfaceData, error)
}
