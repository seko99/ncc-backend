package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_ip_pools_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository IpPools

const (
	IpPoolCreatedEvent    = "IpPoolCreated"
	IpPoolUpdatedEvent    = "IpPoolUpdated"
	IpPoolDeletedEvent    = "IpPoolDeleted"
	IpPoolAllDeletedEvent = "IpPoolAllDeleted"
)

type IpPools interface {
	Create(data []models.IpPoolData) error
	Upsert(data models.IpPoolData) error
	Update(data models.IpPoolData) error
	Delete(id string) error
	DeleteAll() error
	Get() ([]models.IpPoolData, error)
}
