package exporter

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/domain"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
	"time"
)

func (ths *Exporter) exportCustomerServices() ([]exporter.ExportData, error) {
	ths.log.Info("Getting customers...")

	customers, err := ths.customerRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get customers: %w", err)
	}

	var records []exporter.ExportData
	var sormCustomerServicesDiff []models.SormCustomerServiceData

	currentSormCustomerServices, err := ths.sormCustomerServicesRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get SORM customer services: %w", err)
	}

	ths.log.Info("Making customer services diff...")

	for _, c := range customers {
		var record domain.SormCustomerServicesRecord
		record = ths.fromCustomerServiceData(c)
		if !ths.findInSormCustomerServices(currentSormCustomerServices, record.ToSormCustomerServiceData()) {
			records = append(records, record)
			sormCustomerServicesDiff = append(sormCustomerServicesDiff, record.ToSormCustomerServiceData())
		}
	}

	var recordsExported = 0
	var errorsDetected = 0

	ths.log.Info("sormCustomerServicesDiff=%d", len(sormCustomerServicesDiff))

	if len(sormCustomerServicesDiff) > 0 {
		err = ths.sormCustomerServicesErrorsRepo.DeleteAll()
		if err != nil {
			return nil, fmt.Errorf("can't cleanup error customers: %w", err)
		}

		exportTime := time.Now()

		recordsExported = len(records)

		ths.log.Info("Exporting customer services...")

		err = ths.exportData(records)
		if err != nil {
			return nil, fmt.Errorf("can't export customers: %w", err)
		}

		time.Sleep(10 * time.Second)

		ths.log.Info("Checking for customer services errors...")

		exportErrors, err := ths.getErrors(exportTime, "abonents", "services", domain.SormCustomerServicesRecord{})
		if err != nil {
			return nil, err
		}

		errorsDetected = len(exportErrors)

		var errorCustomerServices []models.SormCustomerServicesErrorsData
		for _, ec := range exportErrors {
			record, ok := ec.(domain.SormCustomerServicesRecord)
			if ok {
				errorCustomerServices = append(errorCustomerServices, ths.toSormCustomerServicesError(record, "unknown reason", customers))
			} else {
				ths.log.Error("Can't cast to record: %+v", ec)
			}
		}

		var sormCustomerServicesToUpsert []models.SormCustomerServiceData

		if len(errorCustomerServices) > 0 {
			ths.log.Warn("Error customers detected: %d", len(errorCustomerServices))

			err = ths.sormCustomerServicesErrorsRepo.Create(errorCustomerServices)
			if err != nil {
				return nil, fmt.Errorf("can't create error customer services: %w", err)
			}
		}

		for _, sc := range sormCustomerServicesDiff {
			if !ths.findByLoginInSormCustomerServicesErrors(errorCustomerServices, sc) {
				sormCustomerServicesToUpsert = append(sormCustomerServicesToUpsert, sc)
			}
		}

		ths.log.Info("sormCustomerServicesToUpsert=%d", len(sormCustomerServicesToUpsert))

		if len(sormCustomerServicesToUpsert) > 0 {
			err = ths.sormCustomerServicesRepo.Upsert(sormCustomerServicesToUpsert)
			if err != nil {
				return nil, fmt.Errorf("can't upsert SORM customer services: %w", err)
			}
		}
	}

	var status = "OK"
	if errorsDetected > 0 {
		status = "ERRORS"
	}

	err = ths.sormExportStatusRepo.Upsert(models.SormExportStatusData{
		FileName:    domain.SormCustomerServicesRecord{}.FileName(),
		ExportCount: recordsExported,
		Errors:      errorsDetected,
		Status:      status,
	})
	if err != nil {
		return nil, fmt.Errorf("can't set export status: %w", err)
	}

	return records, nil
}

func (ths *Exporter) findInSormCustomerServices(sormCustomers []models.SormCustomerServiceData, c models.SormCustomerServiceData) bool {
	for _, sc := range sormCustomers {
		if sc.Hash == c.Hash {
			return true
		}
	}
	return false
}

func (ths *Exporter) findByLoginInSormCustomerServicesErrors(sormCustomers []models.SormCustomerServicesErrorsData, c models.SormCustomerServiceData) bool {
	for _, sc := range sormCustomers {
		if sc.Login == c.Login {
			return true
		}
	}
	return false
}

func (ths *Exporter) toSormCustomerServicesError(r domain.SormCustomerServicesRecord, reason string, customers []models.CustomerData) models.SormCustomerServicesErrorsData {
	for _, c := range customers {
		if c.Login == r.Login {
			return models.SormCustomerServicesErrorsData{
				Login:      r.Login,
				CustomerId: models.NewNullUUID(c.Id),
				Reason:     reason,
			}
		}
	}
	return models.SormCustomerServicesErrorsData{}
}

func (ths *Exporter) fromCustomerServiceData(c models.CustomerData) domain.SormCustomerServicesRecord {
	return domain.SormCustomerServicesRecord{
		OrgUnit:        "1",
		Login:          c.Login,
		ContractNumber: c.Contract.Number,
		ServiceCode:    c.ServiceInternet.Code,
		Start:          c.CreateTs.Format(DateFormat),
		End:            "",
		CustomData:     "",
	}
}
