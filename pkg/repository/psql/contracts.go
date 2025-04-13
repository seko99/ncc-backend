package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type Contracts struct {
	storage *psqlstorage.Storage
}

func (s *Contracts) Create(data models.ContractData) error {
	r := s.storage.GetDB().Model(&models.ContractData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Contracts) Update(data models.ContractData) error {
	r := s.storage.GetDB().Model(&models.ContractData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *Contracts) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.ContractData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Contracts) Get() ([]models.ContractData, error) {
	var data []models.ContractData

	r := s.storage.GetDB().Model(&models.ContractData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewContracts(storage *psqlstorage.Storage) *Contracts {
	return &Contracts{
		storage: storage,
	}
}
