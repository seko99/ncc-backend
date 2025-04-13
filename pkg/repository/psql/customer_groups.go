package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type CustomerGroups struct {
	storage *psqlstorage.Storage
}

func (s *CustomerGroups) Create(data models.CustomerGroupData) error {
	r := s.storage.GetDB().Model(&models.CustomerGroupData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *CustomerGroups) Upsert(data models.CustomerGroupData) error {
	r := s.storage.GetDB().Model(&models.CustomerGroupData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *CustomerGroups) Update(data models.CustomerGroupData) error {
	r := s.storage.GetDB().Model(&models.CustomerGroupData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *CustomerGroups) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.CustomerGroupData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *CustomerGroups) Get() ([]models.CustomerGroupData, error) {
	var data []models.CustomerGroupData

	r := s.storage.GetDB().Model(&models.CustomerGroupData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewCustomerGroups(storage *psqlstorage.Storage) *CustomerGroups {
	return &CustomerGroups{
		storage: storage,
	}
}
