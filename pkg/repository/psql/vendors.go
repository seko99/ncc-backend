package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type Vendors struct {
	storage *psqlstorage.Storage
}

func (s *Vendors) Create(data models.VendorData) error {
	r := s.storage.GetDB().Model(&models.VendorData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *Vendors) Upsert(data models.VendorData) error {
	r := s.storage.GetDB().Model(&models.VendorData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *Vendors) Update(data models.VendorData) error {
	r := s.storage.GetDB().Model(&models.VendorData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *Vendors) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.VendorData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Vendors) Get() ([]models.VendorData, error) {
	var data []models.VendorData

	r := s.storage.GetDB().Model(&models.VendorData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewVendors(storage *psqlstorage.Storage) *Vendors {
	return &Vendors{
		storage: storage,
	}
}
