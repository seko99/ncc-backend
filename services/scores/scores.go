package scores

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	repository2 "code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"fmt"
	"math"
	"time"
)

type Scores struct {
	log          *zero.Logger
	customerRepo repository2.Customers
	paymentsRepo repository2.Payments
	scoresRepo   repository2.Scores
}

func NewScores(
	log *zero.Logger,
	customerRepo repository2.Customers,
	paymentsRepo repository2.Payments,
	scoresRepo repository2.Scores,
) *Scores {
	return &Scores{
		log:          log,
		customerRepo: customerRepo,
		paymentsRepo: paymentsRepo,
		scoresRepo:   scoresRepo,
	}
}

func (s *Scores) Process(dryRun bool) error {

	lastScore, err := s.scoresRepo.GetLastScore()
	if err != nil {
		return fmt.Errorf("can't get last scores: %w", err)
	}

	paymentTypes, err := s.scoresRepo.GetPaymentTypes()
	if err != nil {
		return fmt.Errorf("can't get score payment types: %w", err)
	}

	paymentsToScore, err := s.scoresRepo.GetPaymentsToScore(paymentTypes, lastScore.CreateTs)
	if err != nil {
		return fmt.Errorf("can't get payments to score: %w", err)
	}

	ts := time.Now()
	for _, p := range paymentsToScore {
		pType, err := s.scoresRepo.GetScorePaymentTypeForPayment(p, paymentTypes)
		if err != nil {
			s.log.Error("can't process payment: %v", err)
			continue
		}

		scores := math.Round(math.Floor(p.Amount/pType.PaymentBaseAmount) * float64(pType.ScoresAmount))

		if scores > 0 {
			s.log.Info("Creating score for %s (%d) => (%v %0.2f): %d scores", p.Customer.Login, p.Customer.Scores, p.CreateTs, p.Amount, int(scores))

			if !dryRun {
				err := s.scoresRepo.Create(p, int(scores))
				if err != nil {
					s.log.Error("can't create score log: %v", err)
				}
			}
		}
	}
	s.log.Info("Processed %d payments in %v", len(paymentsToScore), time.Since(ts))

	return nil
}
