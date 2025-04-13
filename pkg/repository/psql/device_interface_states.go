package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type DeviceInterfaceStates struct {
	storage *psqlstorage.Storage
}

func (s *DeviceInterfaceStates) Create(data models.IfaceStateData) error {
	r := s.storage.GetDB().Model(&models.IfaceStateData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *DeviceInterfaceStates) Update(data models.IfaceStateData) error {
	r := s.storage.GetDB().Model(&models.IfaceStateData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *DeviceInterfaceStates) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.IfaceStateData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *DeviceInterfaceStates) Get() ([]models.IfaceStateData, error) {
	var data []models.IfaceStateData

	r := s.storage.GetDB().Model(&models.IfaceStateData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewDeviceInterfaceStates(storage *psqlstorage.Storage) *DeviceInterfaceStates {
	return &DeviceInterfaceStates{
		storage: storage,
	}
}
