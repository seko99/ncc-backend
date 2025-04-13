package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type SormCustomerServicesErrors struct {
	storage *psqlstorage.Storage
}

func (s *SormCustomerServicesErrors) Create(data []models.SormCustomerServicesErrorsData) error {
	r := s.storage.GetDB().Model(&models.SormCustomerServicesErrorsData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormCustomerServicesErrors) Upsert(data []models.SormCustomerServicesErrorsData) error {
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

func (s *SormCustomerServicesErrors) Update(data models.SormCustomerServicesErrorsData) error {
	r := s.storage.GetDB().Model(&models.SormCustomerServicesErrorsData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *SormCustomerServicesErrors) DeleteAll() error {
	r := s.storage.GetDB().Exec("DELETE FROM ?", clause.Table{Name: models.SormCustomerServicesErrorsData{}.TableName()})
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *SormCustomerServicesErrors) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.SormCustomerServicesErrorsData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormCustomerServicesErrors) Get() ([]models.SormCustomerServicesErrorsData, error) {
	var data []models.SormCustomerServicesErrorsData

	r := s.storage.GetDB().Model(&models.SormCustomerServicesErrorsData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("login").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewSormCustomerServicesErrors(storage *psqlstorage.Storage) *SormCustomerServicesErrors {
	return &SormCustomerServicesErrors{
		storage: storage,
	}
}
