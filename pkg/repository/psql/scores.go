package psql

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
	"time"
)

type Scores struct {
	storage *psqlstorage.Storage
}

func (s *Scores) Create(payment models2.PaymentData, scores int) error {
	scorePayment := &models2.ScoreLogData{
		CustomerId:   payment.CustomerId,
		Scores:       scores,
		PaymentId:    models2.NewNullUUID(payment.Id),
		ScoresBefore: payment.Customer.Scores,
	}

	tx := s.storage.GetDB().Begin()

	r := tx.Model(&models2.ScoreLogData{}).
		Create(&scorePayment)
	if r.Error != nil {
		tx.Rollback()
		return fmt.Errorf("can't create score log: %w", r.Error)
	}

	r = tx.Model(&models2.CustomerData{}).
		Where("id = ?", payment.CustomerId).
		Update("scores", payment.Customer.Scores+scores)
	if r.Error != nil {
		tx.Rollback()
		return fmt.Errorf("can't update customer scores: %w", r.Error)
	}

	tx.Commit()

	return nil
}

func (s *Scores) GetLastScore() (*models2.ScoreLogData, error) {
	lastScore := &models2.ScoreLogData{}

	r := s.storage.GetDB().Model(&models2.ScoreLogData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("create_ts desc").
		Limit(1).
		First(&lastScore)
	if r.Error != nil {
		return nil, r.Error
	}

	return lastScore, nil
}

func (s *Scores) GetPaymentTypes() ([]models2.ScorePaymentTypes, error) {
	var types []models2.ScorePaymentTypes

	r := s.storage.GetDB().Model(&models2.ScorePaymentTypes{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Find(&types)
	if r.Error != nil {
		return nil, r.Error
	}

	return types, nil
}

func (s *Scores) GetScorePaymentTypeForPayment(payment models2.PaymentData, types []models2.ScorePaymentTypes) (*models2.ScorePaymentTypes, error) {
	for _, t := range types {
		if t.PaymentTypeId == payment.PaymentTypeId {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("score payment type not found: %s", payment.PaymentTypeId.UUID.String())
}

func (s *Scores) GetPaymentsToScore(paymentTypes []models2.ScorePaymentTypes, ts time.Time) ([]models2.PaymentData, error) {
	var payments []models2.PaymentData

	var types []uuid.NullUUID

	for _, t := range paymentTypes {
		types = append(types, t.PaymentTypeId)
	}

	r := s.storage.GetDB().Model(&models2.PaymentData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("date > ?", ts).
		Where("payment_type_id IN ?", types).
		Find(&payments)
	if r.Error != nil {
		return nil, r.Error
	}

	return payments, nil
}

func (s *Scores) Get(period repository.TimePeriod) ([]models2.ScoreLogData, error) {
	var scores []models2.ScoreLogData

	r := s.storage.GetDB().Model(&models2.ScoreLogData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where(repository.PeriodClause("create_ts", period)).
		Order("create_ts").
		Find(&scores)
	if r.Error != nil {
		return nil, r.Error
	}

	return scores, nil
}

func NewScores(storage *psqlstorage.Storage) *Scores {
	return &Scores{
		storage: storage,
	}
}
