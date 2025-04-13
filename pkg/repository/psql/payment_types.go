package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type PaymentTypes struct {
	storage *psqlstorage.Storage
}

func (s *PaymentTypes) Create(data models.PaymentTypeData) error {
	r := s.storage.GetDB().Model(&models.PaymentTypeData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *PaymentTypes) Upsert(data models.PaymentTypeData) error {
	r := s.storage.GetDB().Model(&models.PaymentTypeData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *PaymentTypes) Update(data models.PaymentTypeData) error {
	r := s.storage.GetDB().Model(&models.PaymentTypeData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *PaymentTypes) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.PaymentTypeData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *PaymentTypes) Get() ([]models.PaymentTypeData, error) {
	var data []models.PaymentTypeData

	r := s.storage.GetDB().Model(&models.PaymentTypeData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewPaymentTypes(storage *psqlstorage.Storage) *PaymentTypes {
	return &PaymentTypes{
		storage: storage,
	}
}
