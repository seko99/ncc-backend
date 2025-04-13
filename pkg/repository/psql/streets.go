package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type Streets struct {
	storage *psqlstorage.Storage
}

func (s *Streets) Create(data models.StreetData) error {
	r := s.storage.GetDB().Model(&models.StreetData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Streets) Upsert(data models.StreetData) error {
	r := s.storage.GetDB().Model(&models.StreetData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Streets) Update(data models.StreetData) error {
	r := s.storage.GetDB().Model(&models.StreetData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *Streets) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.StreetData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Streets) Get() ([]models.StreetData, error) {
	var data []models.StreetData

	r := s.storage.GetDB().Model(&models.StreetData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewStreets(storage *psqlstorage.Storage) *Streets {
	return &Streets{
		storage: storage,
	}
}
