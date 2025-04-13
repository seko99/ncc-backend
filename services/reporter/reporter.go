package reporter

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	repository2 "code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/services/interfaces"
	"fmt"
	"time"
)

type Reporter struct {
	log          logger.Logger
	cfg          config.ReporterConfig
	telegram     interfaces.Telegram
	feeRepo      repository2.Fees
	customerRepo repository2.Customers
	sessionRepo  repository2.Sessions
	snapshotRepo repository2.Snapshots
	paymentsRepo repository2.Payments
	scoresRepo   repository2.Scores
}

func (s *Reporter) Run() error {
	start := time.Now().Add(-time.Hour * 24)

	ts := time.Now()
	fees, err := s.feeRepo.Get(repository2.TimePeriod{
		Start: start,
	})
	if err != nil {
		return fmt.Errorf("can't get fees: %w", err)
	}

	snapshots, err := s.snapshotRepo.Get(repository2.TimePeriod{
		Start: start,
	})
	if err != nil {
		s.log.Error("Can't get snapshots: %v", err)
	}

	sessions, err := s.sessionRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get sessions: %w", err)
	}

	payments, err := s.paymentsRepo.GetPayments(repository2.TimePeriod{
		Start: start,
	})
	if err != nil {
		return fmt.Errorf("can't get payments: %w", err)
	}

	scores, err := s.scoresRepo.Get(repository2.TimePeriod{
		Start: start,
	})

	msg := fmt.Sprintf(`За сутки:
Снятий: %d
Платежей: %d
Начислений баллов: %d
Snapshots: %d
-------------
Активных сессий: %d
-------------
Сгенерировано за %v
`, len(fees), len(payments), len(scores), len(snapshots), len(sessions), time.Since(ts))

	err = s.telegram.Send(s.cfg.Telegram.ChatID, msg)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}
func NewReporter(
	log logger.Logger,
	cfg config.ReporterConfig,
	telegram interfaces.Telegram,
	feeRepo repository2.Fees,
	customerRepo repository2.Customers,
	sessionRepo repository2.Sessions,
	snapshotRepo repository2.Snapshots,
	paymentsRepo repository2.Payments,
	scoresRepo repository2.Scores,
) *Reporter {
	return &Reporter{
		log:          log,
		cfg:          cfg,
		telegram:     telegram,
		feeRepo:      feeRepo,
		customerRepo: customerRepo,
		sessionRepo:  sessionRepo,
		snapshotRepo: snapshotRepo,
		paymentsRepo: paymentsRepo,
		scoresRepo:   scoresRepo,
	}
}
