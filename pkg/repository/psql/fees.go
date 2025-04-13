package psql

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type Fees struct {
	storage *psqlstorage.Storage
}

func (s *Fees) Create(fee models2.FeeLogData) error {
	r := s.storage.GetDB().Model(&models2.FeeLogData{}).Create(&fee)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Fees) CreateBatch(fees []models2.FeeLogData) error {
	r := s.storage.GetDB().Model(&models2.FeeLogData{}).Create(&fees)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Fees) Update(fee models2.FeeLogData) error {
	r := s.storage.GetDB().Model(&models2.FeeLogData{}).
		Where("id = @id", sql.Named("id", fee.Id)).
		Updates(fee)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *Fees) GetByCustomer(id string, period repository.TimePeriod, limit ...int) ([]models2.FeeLogData, error) {
	var fees []models2.FeeLogData

	periodClause := repository.PeriodClause("fee_timestamp", period)
	r := s.storage.GetDB().Model(&models2.FeeLogData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("customer_id = @id", sql.Named("id", id)).
		Where(periodClause)

	if len(limit) > 0 {
		r = r.Limit(limit[0]).Order("fee_timestamp DESC")
	} else {
		r = r.Order("fee_timestamp")
	}

	r = r.Find(&fees)
	if r.Error != nil {
		return nil, r.Error
	}

	return fees, nil
}

func (s *Fees) DeleteByCustomer(id string, period repository.TimePeriod) error {
	periodClause := repository.PeriodClause("fee_timestamp", period)
	r := s.storage.GetDB().Model(&models2.FeeLogData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("customer_id = @id", sql.Named("id", id)).
		Where(periodClause).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Fees) Get(period repository.TimePeriod) ([]models2.FeeLogData, error) {
	var fees []models2.FeeLogData

	r := s.storage.GetDB().Model(&models2.FeeLogData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where(repository.PeriodClause("fee_timestamp", period)).
		Order("fee_timestamp").
		Find(&fees)
	if r.Error != nil {
		return nil, r.Error
	}

	return fees, nil
}

func (s *Fees) GetProcessed(t ...time.Time) ([]models2.CustomerData, error) {
	var customers []models2.CustomerData
	var fees []models2.FeeLogData

	y, m, d := time.Now().Date()
	ds := time.Date(y, m, d, 0, 0, 0, 0, time.Now().Location())
	de := time.Date(y, m, d, 23, 59, 59, 0, time.Now().Location())

	if len(t) > 0 {
		y, m, d := t[0].Date()
		ds = time.Date(y, m, d, 0, 0, 0, 0, t[0].Location())
		de = time.Date(y, m, d, 23, 59, 59, 0, t[0].Location())
	}

	r := s.storage.GetDB().Model(&models2.FeeLogData{}).
		Preload(clause.Associations).
		Where("create_ts >= @day_start", sql.Named("day_start", ds)).
		Where("create_ts <= @day_end", sql.Named("day_end", de)).
		Order("fee_timestamp").
		Find(&fees)
	if r.Error != nil {
		return nil, r.Error
	}

	for _, f := range fees {
		customers = append(customers, f.Customer)
	}

	return customers, nil
}

func (s *Fees) GetProcessedMap(t ...time.Time) (map[string]models2.CustomerData, error) {
	processed, err := s.GetProcessed(t...)
	if err != nil {
		return nil, err
	}

	processedMap := make(map[string]models2.CustomerData, len(processed))
	for _, p := range processed {
		processedMap[p.Login] = p
	}

	return processedMap, nil
}

func NewFees(storage *psqlstorage.Storage) *Fees {
	return &Fees{
		storage: storage,
	}
}
