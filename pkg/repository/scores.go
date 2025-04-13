package repository

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"time"
)

//go:generate mockgen -destination=mocks/mock_scores_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository Scores

type Scores interface {
	Create(payment models2.PaymentData, scores int) error
	GetLastScore() (*models2.ScoreLogData, error)
	GetPaymentTypes() ([]models2.ScorePaymentTypes, error)
	GetScorePaymentTypeForPayment(payment models2.PaymentData, types []models2.ScorePaymentTypes) (*models2.ScorePaymentTypes, error)
	GetPaymentsToScore(paymentTypes []models2.ScorePaymentTypes, ts time.Time) ([]models2.PaymentData, error)
	Get(period TimePeriod) ([]models2.ScoreLogData, error)
}
