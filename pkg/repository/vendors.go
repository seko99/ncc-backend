package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_vendors_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository Vendors

type Vendors interface {
	Create(data models.VendorData) error
	Upsert(data models.VendorData) error
	Update(data models.VendorData) error
	Delete(id string) error
	Get() ([]models.VendorData, error)
}
