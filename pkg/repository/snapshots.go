package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_snapshots_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository Snapshots

type Snapshots interface {
	Create(snapshot *models.SnapshotData) error
	Get(period TimePeriod) ([]models.SnapshotData, error)
}
