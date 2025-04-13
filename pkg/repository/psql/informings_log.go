package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"gorm.io/gorm/clause"
)

type InformingLog struct {
	storage *psqlstorage.Storage
}

func (s *InformingLog) Get() ([]models.InformingLogData, error) {
	var data []models.InformingLogData
	r := s.storage.GetDB().Model(models.InformingLogData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("create_ts").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}
	return data, nil
}

func (s *InformingLog) Create(data []models.InformingLogData) error {
	r := s.storage.GetDB().Model(&models.InformingLogData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *InformingLog) SetStatus(data models.InformingLogData, status int) error {
	r := s.storage.GetDB().Model(&models.InformingLogData{}).
		Where("id = ?", data.Id).
		Update("status", status)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func NewInformingLog(storage *psqlstorage.Storage) *InformingLog {
	return &InformingLog{
		storage: storage,
	}
}
