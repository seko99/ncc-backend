package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type Devices struct {
	storage *psqlstorage.Storage
}

func (s *Devices) Create(fee models.DeviceData) error {
	r := s.storage.GetDB().Model(&models.DeviceData{}).Create(&fee)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Devices) Update(fee models.DeviceData) error {
	r := s.storage.GetDB().Model(&models.DeviceData{}).
		Where("id = @id", sql.Named("id", fee.Id)).
		Updates(fee)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *Devices) Get() ([]models.DeviceData, error) {
	var data []models.DeviceData

	r := s.storage.GetDB().Model(&models.DeviceData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func (s *Devices) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.DeviceData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func NewDevices(storage *psqlstorage.Storage) *Devices {
	return &Devices{
		storage: storage,
	}
}
