package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_radius_attributes_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository RadiusAttributes

type RadiusAttributes interface {
	Create(data models.RadiusAttributeData) error
	Upsert(data models.RadiusAttributeData) error
	Update(data models.RadiusAttributeData) error
	Delete(id string) error
	Get() ([]models.RadiusAttributeData, error)
}
