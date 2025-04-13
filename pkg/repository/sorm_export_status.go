package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_sorm_export_status_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository SormExportStatus

type SormExportStatus interface {
	Create(data models.SormExportStatusData) error
	Upsert(data models.SormExportStatusData) error
	Update(data models.SormExportStatusData) error
	Delete(id string) error
	Get() ([]models.SormExportStatusData, error)
	GetByFileName(fileName string) (models.SormExportStatusData, error)
}
