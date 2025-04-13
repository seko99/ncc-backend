package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type DeviceStates struct {
	storage *psqlstorage.Storage
}

func (s *DeviceStates) Create(data models.DeviceStateData) error {
	r := s.storage.GetDB().Model(&models.DeviceStateData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *DeviceStates) Update(data models.DeviceStateData) error {
	r := s.storage.GetDB().Model(&models.DeviceStateData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *DeviceStates) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.DeviceStateData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *DeviceStates) Get() ([]models.DeviceStateData, error) {
	var data []models.DeviceStateData

	r := s.storage.GetDB().Model(&models.DeviceStateData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewDeviceStates(storage *psqlstorage.Storage) *DeviceStates {
	return &DeviceStates{
		storage: storage,
	}
}
