package cmd

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	psql2 "code.evixo.ru/ncc/ncc-backend/pkg/repository/psql"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	s3storage "code.evixo.ru/ncc/ncc-backend/pkg/storage/s3"
	"code.evixo.ru/ncc/ncc-backend/services/exporter"
	exporter2 "code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var exporterCmd = &cobra.Command{
	Use:   "exporter",
	Short: "Exporter",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.NewConfig()
		if err != nil {
			panic(err)
		}

		debugLevel := zerolog.InfoLevel

		if *debugMode {
			debugLevel = zerolog.DebugLevel
		}
		log := zero.NewLogger(debugLevel)

		log.Info("Starting Exporter")

		storage := psqlstorage.NewStorage(cfg, log, psqlstorage.WithAppName("ncc-exporter"))
		err = storage.Connect()
		if err != nil {
			log.Error("can't connect to storage: %v", err)
			return
		}

		customerRepo := psql2.NewCustomers(storage, nil)
		paymentsRepo := psql2.NewPayments(storage)
		paymentTypesRepo := psql2.NewPaymentTypes(storage)
		serviceInternetRepo := psql2.NewServiceInternet(storage, nil)
		documentTypesRepo := psql2.NewDocumentTypes(storage)
		ipNumberingRepo := psql2.NewSormIpNumbering(storage)
		gatewayRepo := psql2.NewSormGateway(storage)
		sormCustomersRepo := psql2.NewSormCustomers(storage)
		sormCustomersErrorsRepo := psql2.NewSormCustomersErrors(storage)
		sormCustomerServicesRepo := psql2.NewSormCustomerServices(storage)
		sormCustomerServicesErrorsRepo := psql2.NewSormCustomerServicesErrors(storage)
		sormExportStatusRepo := psql2.NewSormExportStatus(storage)

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
				log.Error("Can't create export SSH writer: %v", err)
				return
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
				log.Error("Can't create export FTP writer: %v", err)
				return
			}
		}

		s3storage := s3storage.NewS3(cfg, log)
		err = s3storage.Connect()
		if err != nil {
			log.Error("can't connect to S3 storage: %v", err)
			return
		}

		backupWriter, err := exporter.NewS3Writer(s3storage, "backup_writer", true)

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
			exporter.WithBackupEnabled(backupWriter),
		)
		err = exporterService.Run()
		if err != nil {
			log.Error("can't run exporter: %v", err)
			return
		}
	},
}
