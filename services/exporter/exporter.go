package exporter

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	repository2 "code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
	"time"
)

const (
	DateFormat = "02.01.2006 15:04:05"
)

type Exporter struct {
	log                            logger.Logger
	writer                         exporter.ExportWriter
	backupWriter                   exporter.ExportWriter
	backupEnabled                  bool
	customerRepo                   repository2.Customers
	paymentsRepo                   repository2.Payments
	paymentTypesRepo               repository2.PaymentTypes
	serviceInternetRepo            repository2.ServiceInternet
	documentTypesRepo              repository2.DocumentTypes
	ipNumberingRepo                repository2.SormIpNumbering
	gatewayRepo                    repository2.SormGateway
	sormCustomersRepo              repository2.SormCustomers
	sormCustomersErrorsRepo        repository2.SormCustomersErrors
	sormCustomerServicesRepo       repository2.SormCustomerServices
	sormCustomerServicesErrorsRepo repository2.SormCustomerServicesErrors
	sormExportStatusRepo           repository2.SormExportStatus
}

type ExporterOption func(exporter *Exporter)

func (ths *Exporter) Run() error {
	var err error

	_, err = ths.exportCustomers()
	if err != nil {
		ths.log.Error("Can't export customers: %v", err)
	}

	_, err = ths.exportCustomerServices()
	if err != nil {
		ths.log.Error("Can't export customer services: %v", err)
	}

	_, err = ths.exportPayments()
	if err != nil {
		ths.log.Error("Can't export payments: %v", err)
	}

	_, err = ths.exportServices()
	if err != nil {
		return fmt.Errorf("can't prepare service list: %w", err)
	}

	_, err = ths.exportPaymentTypes()
	if err != nil {
		return fmt.Errorf("can't prepare pay types list: %w", err)
	}

	_, err = ths.exportDocTypes()
	if err != nil {
		return fmt.Errorf("can't prepare document types list: %w", err)
	}

	_, err = ths.exportNumberingPlan()
	if err != nil {
		return fmt.Errorf("can't prepare IP numbering plan list: %w", err)
	}

	_, err = ths.exportGates()
	if err != nil {
		return fmt.Errorf("can't prepare gates list: %w", err)
	}

	return nil
}

func (ths *Exporter) ScheduledRun() error {
	return ths.Run()
}

func (ths *Exporter) exportData(data []exporter.ExportData, withHeader ...bool) error {
	ths.log.Info("Writing to main...")
	err := ths.writer.Write(data, withHeader...)
	if err != nil {
		ths.log.Error("Can't export data: %v", err)
		//return fmt.Errorf("can't export data: %w", err)
	}

	if ths.backupEnabled {
		ths.log.Info("Writing to backup...")
		err = ths.backupWriter.Write(data, withHeader...)
		if err != nil {
			ths.log.Error("Can't export to backup: %v", err)
		}
	}

	return nil
}

func (ths *Exporter) getErrors(exportTime time.Time, path, errorFileName string, d exporter.ExportData) ([]exporter.ExportData, error) {
	exportErrors, err := ths.writer.GetErrors(exportTime, path, errorFileName, d)
	if err != nil {
		return nil, fmt.Errorf("can't get export errors: %w", err)
	}
	return exportErrors, nil
}

func WithBackupEnabled(writer exporter.ExportWriter) ExporterOption {
	return func(exporter *Exporter) {
		exporter.backupWriter = writer
		exporter.backupEnabled = true
	}
}

func NewExporter(
	log logger.Logger,
	writer exporter.ExportWriter,
	customerRepo repository2.Customers,
	paymentsRepo repository2.Payments,
	paymentTypesRepo repository2.PaymentTypes,
	serviceInternetRepo repository2.ServiceInternet,
	documentTypesRepo repository2.DocumentTypes,
	ipNumberingRepo repository2.SormIpNumbering,
	gatewayRepo repository2.SormGateway,
	sormCustomersRepo repository2.SormCustomers,
	sormCustomersErrorsRepo repository2.SormCustomersErrors,
	sormCustomerServicesRepo repository2.SormCustomerServices,
	sormCustomerServicesErrorsRepo repository2.SormCustomerServicesErrors,
	sormExportStatusRepo repository2.SormExportStatus,
	opts ...ExporterOption,
) *Exporter {
	exp := &Exporter{
		log:                            log,
		writer:                         writer,
		customerRepo:                   customerRepo,
		paymentsRepo:                   paymentsRepo,
		paymentTypesRepo:               paymentTypesRepo,
		serviceInternetRepo:            serviceInternetRepo,
		documentTypesRepo:              documentTypesRepo,
		ipNumberingRepo:                ipNumberingRepo,
		gatewayRepo:                    gatewayRepo,
		sormCustomersRepo:              sormCustomersRepo,
		sormCustomersErrorsRepo:        sormCustomersErrorsRepo,
		sormCustomerServicesRepo:       sormCustomerServicesRepo,
		sormCustomerServicesErrorsRepo: sormCustomerServicesErrorsRepo,
		sormExportStatusRepo:           sormExportStatusRepo,
	}

	for _, opt := range opts {
		opt(exp)
	}

	return exp
}
