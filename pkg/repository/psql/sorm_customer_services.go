package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type SormCustomerServices struct {
	storage *psqlstorage.Storage
}

func (s *SormCustomerServices) Create(data []models.SormCustomerServiceData) error {
	r := s.storage.GetDB().Model(&models.SormCustomerServiceData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormCustomerServices) Upsert(data []models.SormCustomerServiceData) error {
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

func (s *SormCustomerServices) Update(data models.SormCustomerServiceData) error {
	r := s.storage.GetDB().Model(&models.SormCustomerServiceData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *SormCustomerServices) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.SormCustomerServiceData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormCustomerServices) Get() ([]models.SormCustomerServiceData, error) {
	var data []models.SormCustomerServiceData

	r := s.storage.GetDB().Model(&models.SormCustomerServiceData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("login").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewSormCustomerServices(storage *psqlstorage.Storage) *SormCustomerServices {
	return &SormCustomerServices{
		storage: storage,
	}
}
