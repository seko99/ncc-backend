package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_document_types_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository DocumentTypes

type DocumentTypes interface {
	Create(data models.DocumentTypeData) error
	Upsert(data models.DocumentTypeData) error
	Update(data models.DocumentTypeData) error
	Delete(id string) error
	Get() ([]models.DocumentTypeData, error)
}
