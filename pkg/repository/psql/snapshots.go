package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"gorm.io/gorm/clause"
)

type Snapshots struct {
	storage *psqlstorage.Storage
}

func (s *Snapshots) Create(snapshot *models.SnapshotData) error {
	r := s.storage.GetDB().Model(&models.SnapshotData{}).Create(snapshot)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *Snapshots) Get(period repository.TimePeriod) ([]models.SnapshotData, error) {
	var snapshots []models.SnapshotData

	r := s.storage.GetDB().Model(&models.SnapshotData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where(repository.PeriodClause("create_ts", period)).
		Order("create_ts").
		Find(&snapshots)
	if r.Error != nil {
		return nil, r.Error
	}

	return snapshots, nil
}

func NewSnapshots(storage *psqlstorage.Storage) *Snapshots {
	return &Snapshots{
		storage: storage,
	}
}
