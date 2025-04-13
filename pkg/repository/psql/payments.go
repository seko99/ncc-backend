package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
)

type Payments struct {
	storage *psqlstorage.Storage
}

func (s *Payments) GetPaymentsByCustomer(id string, period repository.TimePeriod) ([]models.PaymentData, error) {
	var payments []models.PaymentData

	r := s.storage.GetDB().Model(&models.PaymentData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("customer_id = @id", sql.Named("id", id)).
		Where(repository.PeriodClause("date", period)).
		Order("date DESC").
		Find(&payments)

	if r.Error != nil {
		return nil, r.Error
	}

	return payments, nil
}

func (s *Payments) Create(data models.PaymentData) error {
	r := s.storage.GetDB().Model(&models.PaymentData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Payments) GetPayments(period repository.TimePeriod) ([]models.PaymentData, error) {
	var payments []models.PaymentData

	r := s.storage.GetDB().Model(&models.PaymentData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where(repository.PeriodClause("date", period)).
		Order("date DESC").
		Find(&payments)

	if r.Error != nil {
		return nil, r.Error
	}

	return payments, nil
}

func NewPayments(storage *psqlstorage.Storage) *Payments {
	return &Payments{
		storage: storage,
	}
}
