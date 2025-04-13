package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"time"
)

//go:generate mockgen -destination=mocks/mock_informings_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository Informings

type Informings interface {
	Get() ([]models.InformingData, error)
	GetEnabled() ([]models.InformingData, error)
	Create(data models.InformingData) error
	SetState(data models.InformingData, state int) error
	SetStart(data models.InformingData, start time.Time) error
}
