package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type PaymentSystems struct {
	storage *psqlstorage.Storage
}

func (s *PaymentSystems) Create(data models.PaymentSystemData) error {
	r := s.storage.GetDB().Model(&models.PaymentSystemData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *PaymentSystems) Upsert(data models.PaymentSystemData) error {
	r := s.storage.GetDB().Model(&models.PaymentSystemData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *PaymentSystems) Update(data models.PaymentSystemData) error {
	r := s.storage.GetDB().Model(&models.PaymentSystemData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *PaymentSystems) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.PaymentSystemData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *PaymentSystems) Get() ([]models.PaymentSystemData, error) {
	var data []models.PaymentSystemData

	r := s.storage.GetDB().Model(&models.PaymentSystemData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewPaymentSystems(storage *psqlstorage.Storage) *PaymentSystems {
	return &PaymentSystems{
		storage: storage,
	}
}
