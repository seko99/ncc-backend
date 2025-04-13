package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_dhcp_leases_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository DhcpLeases

const (
	LeaseCreatedEvent    = "leaseCreated"
	LeaseUpdatedEvent    = "leaseUpdated"
	LeaseDeletedEvent    = "leaseDeleted"
	LeaseAllDeletedEvent = "leaseAllDeleted"
)

type DhcpLeases interface {
	Get() ([]models.LeaseData, error)
	GetByIP(ip string) (models.LeaseData, error)

	Create(data models.LeaseData) error
	Update(data models.LeaseData) error

	DeleteAll() error
}
