package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type DeviceInterfaces struct {
	storage *psqlstorage.Storage
}

func (s *DeviceInterfaces) Create(data models.IfaceData) error {
	r := s.storage.GetDB().Model(&models.IfaceData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *DeviceInterfaces) Update(data models.IfaceData) error {
	r := s.storage.GetDB().Model(&models.IfaceData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *DeviceInterfaces) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.IfaceData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *DeviceInterfaces) Get() ([]models.IfaceData, error) {
	var data []models.IfaceData

	r := s.storage.GetDB().Model(&models.IfaceData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("port").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func (s *DeviceInterfaces) GetByDeviceIdAndPort(id string, port int) (models.IfaceData, error) {
	var data models.IfaceData

	r := s.storage.GetDB().Model(&models.IfaceData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("device = @device_id", sql.Named("device_id", id)).
		Where("port = @port", sql.Named("port", port)).
		Order("port").
		Find(&data)
	if r.Error != nil {
		return models.IfaceData{}, r.Error
	}

	return data, nil
}

func NewDeviceInterfaces(storage *psqlstorage.Storage) *DeviceInterfaces {
	return &DeviceInterfaces{
		storage: storage,
	}
}
