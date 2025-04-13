package scheduler

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository/psql"
	scheduler2 "code.evixo.ru/ncc/ncc-backend/pkg/scheduler"
	psqlstorage "code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"code.evixo.ru/ncc/ncc-backend/services/exporter"
	exporter2 "code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
	"time"
)

func RegisterExporter(cfg *config.Config, log logger.Logger, scheduler *Scheduler) error {
	storage := psqlstorage.NewStorage(cfg, log)
	err := storage.Connect()
	if err != nil {
		return fmt.Errorf("can't connect to storage: %v", err)
	}

	customerRepo := psql.NewCustomers(storage, nil)
	paymentsRepo := psql.NewPayments(storage)
	paymentTypesRepo := psql.NewPaymentTypes(storage)
	serviceInternetRepo := psql.NewServiceInternet(storage, nil)
	documentTypesRepo := psql.NewDocumentTypes(storage)
	ipNumberingRepo := psql.NewSormIpNumbering(storage)
	gatewayRepo := psql.NewSormGateway(storage)
	sormCustomersRepo := psql.NewSormCustomers(storage)
	sormCustomersErrorsRepo := psql.NewSormCustomersErrors(storage)
	sormCustomerServicesRepo := psql.NewSormCustomerServices(storage)
	sormCustomerServicesErrorsRepo := psql.NewSormCustomerServicesErrors(storage)
	sormExportStatusRepo := psql.NewSormExportStatus(storage)

	var exportWriter exporter2.ExportWriter

	switch cfg.Exporter.Type {
	case "ssh":
		exportWriter, err = exporter.NewSshWriter(
			cfg.Exporter.Host,
			cfg.Exporter.Username,
			cfg.Exporter.Key,
			cfg.Exporter.Path,
		)
		if err != nil {
			return fmt.Errorf("can't create export SSH writer: %v", err)
		}
	case "ftp":
		exportWriter, err = exporter.NewFtpWriter(
			cfg.Exporter.Host,
			cfg.Exporter.Username,
			cfg.Exporter.Password,
			cfg.Exporter.BadsHost,
			cfg.Exporter.BadsUsername,
			cfg.Exporter.BadsPassword,
			cfg.Exporter.Path,
			true,
		)
		if err != nil {
			return fmt.Errorf("can't create export FTP writer: %v", err)
		}
	}

	exporterService := exporter.NewExporter(
		log,
		exportWriter,
		customerRepo,
		paymentsRepo,
		paymentTypesRepo,
		serviceInternetRepo,
		documentTypesRepo,
		ipNumberingRepo,
		gatewayRepo,
		sormCustomersRepo,
		sormCustomersErrorsRepo,
		sormCustomerServicesRepo,
		sormCustomerServicesErrorsRepo,
		sormExportStatusRepo,
	)

	scheduler.RegisterTask(scheduler2.Task{
		Name:      "exporter",
		IsEnabled: true,
		Task:      exporterService,
		Schedule: scheduler2.Schedule{
			Every: time.Hour,
		},
	})

	return nil
}
