package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_informing_log_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository InformingLog

type InformingLog interface {
	Get() ([]models.InformingLogData, error)
	Create(data []models.InformingLogData) error
	SetStatus(data models.InformingLogData, status int) error
}
