package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type DocumentTypes struct {
	storage *psqlstorage.Storage
}

func (s *DocumentTypes) Create(data models.DocumentTypeData) error {
	r := s.storage.GetDB().Model(&models.DocumentTypeData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *DocumentTypes) Upsert(data models.DocumentTypeData) error {
	r := s.storage.GetDB().Model(&models.DocumentTypeData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *DocumentTypes) Update(data models.DocumentTypeData) error {
	r := s.storage.GetDB().Model(&models.DocumentTypeData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *DocumentTypes) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.DocumentTypeData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *DocumentTypes) Get() ([]models.DocumentTypeData, error) {
	var data []models.DocumentTypeData

	r := s.storage.GetDB().Model(&models.DocumentTypeData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewDocumentTypes(storage *psqlstorage.Storage) *DocumentTypes {
	return &DocumentTypes{
		storage: storage,
	}
}
