package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type SormCustomersErrors struct {
	storage *psqlstorage.Storage
}

func (s *SormCustomersErrors) Create(data []models.SormCustomersErrorsData) error {
	r := s.storage.GetDB().Model(&models.SormCustomersErrorsData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormCustomersErrors) Upsert(data []models.SormCustomersErrorsData) error {
	r := s.storage.GetDB().
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "login"}},
			UpdateAll: true,
		},
			clause.Returning{}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormCustomersErrors) Update(data models.SormCustomersErrorsData) error {
	r := s.storage.GetDB().Model(&models.SormCustomersErrorsData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *SormCustomersErrors) DeleteAll() error {
	r := s.storage.GetDB().Exec("DELETE FROM ?", clause.Table{Name: models.SormCustomersErrorsData{}.TableName()})
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *SormCustomersErrors) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.SormCustomersErrorsData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormCustomersErrors) Get() ([]models.SormCustomersErrorsData, error) {
	var data []models.SormCustomersErrorsData

	r := s.storage.GetDB().Model(&models.SormCustomersErrorsData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("login").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewSormCustomersErrors(storage *psqlstorage.Storage) *SormCustomersErrors {
	return &SormCustomersErrors{
		storage: storage,
	}
}
