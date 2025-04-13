package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_dhcp_bindings_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository DhcpBindings

const (
	BindingCreatedEvent    = "bindingCreated"
	BindingUpdatedEvent    = "bindingUpdated"
	BindingDeletedEvent    = "bindingDeleted"
	BindingAllDeletedEvent = "bindingAllDeleted"
)

type DhcpBindings interface {
	Create(data []models.DhcpBindingData) error
	Update(data models.DhcpBindingData) error
	Delete(id string) error
	DeleteAll() error
	Get() ([]models.DhcpBindingData, error)
}
