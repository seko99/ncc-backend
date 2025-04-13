package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_device_states_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository DeviceStates

type DeviceStates interface {
	Create(data models.DeviceStateData) error
	Update(data models.DeviceStateData) error
	Delete(id string) error
	Get() ([]models.DeviceStateData, error)
}
