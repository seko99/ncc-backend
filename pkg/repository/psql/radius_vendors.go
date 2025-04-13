package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type RadiusVendors struct {
	storage *psqlstorage.Storage
}

func (s *RadiusVendors) Create(data models.RadiusVendorData) error {
	r := s.storage.GetDB().Model(&models.RadiusVendorData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *RadiusVendors) Upsert(data models.RadiusVendorData) error {
	r := s.storage.GetDB().Model(&models.RadiusVendorData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *RadiusVendors) Update(data models.RadiusVendorData) error {
	r := s.storage.GetDB().Model(&models.RadiusVendorData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *RadiusVendors) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.RadiusVendorData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *RadiusVendors) Get() ([]models.RadiusVendorData, error) {
	var data []models.RadiusVendorData

	r := s.storage.GetDB().Model(&models.RadiusVendorData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewRadiusVendors(storage *psqlstorage.Storage) *RadiusVendors {
	return &RadiusVendors{
		storage: storage,
	}
}
