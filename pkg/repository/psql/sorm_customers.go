package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type SormCustomers struct {
	storage *psqlstorage.Storage
}

func (s *SormCustomers) Create(data []models.SormCustomersData) error {
	r := s.storage.GetDB().Model(&models.SormCustomersData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormCustomers) Upsert(data []models.SormCustomersData) error {
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

func (s *SormCustomers) Update(data models.SormCustomersData) error {
	r := s.storage.GetDB().Model(&models.SormCustomersData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *SormCustomers) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.SormCustomersData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormCustomers) Get() ([]models.SormCustomersData, error) {
	var data []models.SormCustomersData

	r := s.storage.GetDB().Model(&models.SormCustomersData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("login").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewSormCustomers(storage *psqlstorage.Storage) *SormCustomers {
	return &SormCustomers{
		storage: storage,
	}
}
