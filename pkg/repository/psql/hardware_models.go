package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type HardwareModels struct {
	storage *psqlstorage.Storage
}

func (s *HardwareModels) Create(data models.HardwareModelData) error {
	r := s.storage.GetDB().Model(&models.HardwareModelData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *HardwareModels) Upsert(data models.HardwareModelData) error {
	r := s.storage.GetDB().Model(&models.HardwareModelData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *HardwareModels) Update(data models.HardwareModelData) error {
	r := s.storage.GetDB().Model(&models.HardwareModelData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *HardwareModels) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.HardwareModelData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *HardwareModels) Get() ([]models.HardwareModelData, error) {
	var data []models.HardwareModelData

	r := s.storage.GetDB().Model(&models.HardwareModelData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewHardwareModels(storage *psqlstorage.Storage) *HardwareModels {
	return &HardwareModels{
		storage: storage,
	}
}
