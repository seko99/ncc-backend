package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_dhcp_pools_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository DhcpPools

const (
	PoolCreatedEvent    = "poolCreated"
	PoolUpdatedEvent    = "poolUpdated"
	PoolDeletedEvent    = "poolDeleted"
	PoolAllDeletedEvent = "poolAllDeleted"
)

type DhcpPools interface {
	Create(data []models.DhcpPoolData) error
	Upsert(data models.DhcpPoolData) error
	Update(data models.DhcpPoolData) error
	Delete(id string) error
	DeleteAll() error
	Get() ([]models.DhcpPoolData, error)
}
