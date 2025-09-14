package exporter

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/domain"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
	"time"
)

func (ths *Exporter) exportPayments() ([]exporter.ExportData, error) {
	var records []exporter.ExportData
	lastExport := time.Now().Add(-4 * 30 * 86400 * time.Second)

	status, err := ths.sormExportStatusRepo.GetByFileName(domain.SormPaymentsRecord{}.FileName())
	if err != nil {
		return nil, fmt.Errorf("can't get export status: %w", err)
	}

	if !status.UpdateTs.IsZero() {
		lastExport = status.CommonData.UpdateTs
	}

	ths.log.Info("Getting payments since %s...", lastExport.Format("2006-01-02 15:04:05"))

	payments, err := ths.paymentsRepo.GetPayments(repository.TimePeriod{
		Start: lastExport,
	})
	if err != nil {
		return nil, fmt.Errorf("can't get payments from %v: %w", lastExport, err)
	}

	var recordsExported = 0

	if len(payments) > 0 {
		for _, p := range payments {
			records = append(records, domain.SormPaymentsRecord{
				OrgUnit:     "1",
				PaymentCode: p.PaymentType.Code,
				Login:       p.Customer.Login,
				Date:        p.Date.Format(DateFormat),
				Amount:      fmt.Sprintf("%0.2f", p.Amount),
			})
		}

		recordsExported = len(records)

		ths.log.Info("Exporting payments...")

		err = ths.exportData(records)
		if err != nil {
			return nil, fmt.Errorf("can't export payments: %w", err)
		}
	}

	err = ths.sormExportStatusRepo.Upsert(models.SormExportStatusData{
		FileName:    domain.SormPaymentsRecord{}.FileName(),
		ExportCount: recordsExported,
		Errors:      0,
		Status:      "OK",
	})
	if err != nil {
		return nil, fmt.Errorf("can't set export status: %w", err)
	}

	return records, nil
}

func (ths *Exporter) exportPaymentTypes() ([]exporter.ExportData, error) {
	status, err := ths.sormExportStatusRepo.GetByFileName(domain.SormPayTypesRecord{}.FileName())
	if err != nil {
		return nil, fmt.Errorf("can't get export status: %w", err)
	}

	if status.Status != "PENDING" {
		return []exporter.ExportData{}, nil
	}

	paymentTypes, err := ths.paymentTypesRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get payment types: %w", err)
	}

	var records []exporter.ExportData

	for _, s := range paymentTypes {
		records = append(records, domain.SormPayTypesRecord{
			OrgUnit:     "1",
			ServiceCode: s.Code,
			Start:       s.CreateTs.Format(DateFormat),
			End:         "",
			Descr:       s.Name,
		})
	}

	err = ths.exportData(records)
	if err != nil {
		return nil, fmt.Errorf("can't export pay types: %w", err)
	}

	err = ths.sormExportStatusRepo.Upsert(models.SormExportStatusData{
		FileName:    domain.SormPayTypesRecord{}.FileName(),
		ExportCount: len(records),
		Errors:      0,
		Status:      "OK",
	})
	if err != nil {
		return nil, fmt.Errorf("can't set export status: %w", err)
	}

	return records, nil
}
