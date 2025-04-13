package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_device_interface_states_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository DeviceInterfaceStates

type DeviceInterfaceStates interface {
	Create(data models.IfaceStateData) error
	Update(data models.IfaceStateData) error
	Delete(id string) error
	Get() ([]models.IfaceStateData, error)
}
