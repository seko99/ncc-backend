package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type Cities struct {
	storage *psqlstorage.Storage
}

func (s *Cities) Create(data models.CityData) error {
	r := s.storage.GetDB().Model(&models.CityData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Cities) Upsert(data models.CityData) error {
	r := s.storage.GetDB().Model(&models.CityData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Cities) Update(data models.CityData) error {
	r := s.storage.GetDB().Model(&models.CityData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *Cities) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.CityData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Cities) Get() ([]models.CityData, error) {
	var data []models.CityData

	r := s.storage.GetDB().Model(&models.CityData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewCities(storage *psqlstorage.Storage) *Cities {
	return &Cities{
		storage: storage,
	}
}
