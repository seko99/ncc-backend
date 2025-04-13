package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_radius_vendors_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository RadiusVendors

type RadiusVendors interface {
	Create(data models.RadiusVendorData) error
	Upsert(data models.RadiusVendorData) error
	Update(data models.RadiusVendorData) error
	Delete(id string) error
	Get() ([]models.RadiusVendorData, error)
}
