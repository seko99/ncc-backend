package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type SormExportStatus struct {
	storage *psqlstorage.Storage
}

func (s *SormExportStatus) Create(data models.SormExportStatusData) error {
	r := s.storage.GetDB().Model(&models.SormExportStatusData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormExportStatus) Upsert(data models.SormExportStatusData) error {
	data.CommonData.UpdateTs = time.Now()
	r := s.storage.GetDB().Model(&models.SormExportStatusData{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "file_name"}},
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormExportStatus) Update(data models.SormExportStatusData) error {
	r := s.storage.GetDB().Model(&models.SormExportStatusData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *SormExportStatus) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.SormExportStatusData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormExportStatus) Get() ([]models.SormExportStatusData, error) {
	var data []models.SormExportStatusData

	r := s.storage.GetDB().Model(&models.SormExportStatusData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("file_name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func (s *SormExportStatus) GetByFileName(fileName string) (models.SormExportStatusData, error) {
	var data models.SormExportStatusData

	r := s.storage.GetDB().Model(&models.SormExportStatusData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("file_name = @file_name", sql.Named("file_name", fileName)).
		Order("file_name").
		Find(&data)
	if r.Error != nil {
		return models.SormExportStatusData{}, r.Error
	}

	return data, nil
}

func NewSormExportStatus(storage *psqlstorage.Storage) *SormExportStatus {
	return &SormExportStatus{
		storage: storage,
	}
}
