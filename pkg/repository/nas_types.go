package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_nas_types_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository NasTypes

type NasTypes interface {
	Create(data models.NasTypeData) error
	Upsert(data models.NasTypeData) error
	Get() ([]models.NasTypeData, error)
}
