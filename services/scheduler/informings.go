package scheduler

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/pkg/providers"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository/psql"
	scheduler2 "code.evixo.ru/ncc/ncc-backend/pkg/scheduler"
	psqlstorage "code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"code.evixo.ru/ncc/ncc-backend/services/informings"
	"fmt"
	"time"
)

func RegisterInformings(cfg *config.Config, log logger.Logger, scheduler *Scheduler) error {
	storage := psqlstorage.NewStorage(cfg, log)
	err := storage.Connect()
	if err != nil {
		return fmt.Errorf("can't connect to storage: %v", err)
	}

	informingsRepo := psql.NewInformings(storage)
	informingLogRepo := psql.NewInformingLog(storage)
	informingsTestCustomersRepo := psql.NewInformingsTestCustomers(storage)
	customerRepo := psql.NewCustomers(storage, nil)
	phoenixSms := providers.NewPhoenixSms(cfg.Informings.SmsProvider, log)
	informingsService := informings.NewInformings(
		log,
		phoenixSms,
		informingsRepo,
		informingLogRepo,
		informingsTestCustomersRepo,
		customerRepo,
	)

	scheduler.RegisterTask(scheduler2.Task{
		Name:      "informings",
		IsEnabled: true,
		Task:      informingsService,
		Schedule: scheduler2.Schedule{
			Every: time.Second * 10,
		},
	})

	return nil
}
