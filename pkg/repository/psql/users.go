package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type Users struct {
	storage *psqlstorage.Storage
}

func (s *Users) Create(data models.SecUserData) error {
	r := s.storage.GetDB().Model(&models.SecUserData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Users) Update(data models.SecUserData) error {
	r := s.storage.GetDB().Model(&models.SecUserData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *Users) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.SecUserData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Users) Get() ([]models.SecUserData, error) {
	var data []models.SecUserData

	r := s.storage.GetDB().Model(&models.SecUserData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func (s *Users) GetGroups() ([]models.SecGroupData, error) {
	var data []models.SecGroupData

	r := s.storage.GetDB().Model(&models.SecGroupData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewUsers(storage *psqlstorage.Storage) *Users {
	return &Users{
		storage: storage,
	}
}
