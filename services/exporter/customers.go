package exporter

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/domain"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
	"strings"
	"time"
)

func (ths *Exporter) exportCustomers() ([]exporter.ExportData, error) {
	ths.log.Info("Getting customers...")

	customers, err := ths.customerRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get customers: %w", err)
	}

	var records []exporter.ExportData
	var sormCustomersDiff []models.SormCustomersData

	currentSormCustomers, err := ths.sormCustomersRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get SORM customers: %w", err)
	}

	ths.log.Info("Making customers diff...")

	for _, c := range customers {
		var record domain.SormCustomersRecord
		record = ths.fromCustomerData(c)
		if !ths.findInSormCustomers(currentSormCustomers, record.ToSormCustomersData()) {
			records = append(records, record)
			sormCustomersDiff = append(sormCustomersDiff, record.ToSormCustomersData())
		}
	}

	var recordsExported = 0
	var errorsDetected = 0

	ths.log.Info("sormCustomersDiff=%d", len(sormCustomersDiff))

	if len(sormCustomersDiff) > 0 {
		err = ths.sormCustomersErrorsRepo.DeleteAll()
		if err != nil {
			return nil, fmt.Errorf("can't cleanup error customers: %w", err)
		}

		exportTime := time.Now()

		recordsExported = len(records)

		ths.log.Info("Exporting customers...")

		err = ths.exportData(records)
		if err != nil {
			return nil, fmt.Errorf("can't export customers: %w", err)
		}

		time.Sleep(30 * time.Second)

		ths.log.Info("Checking for customers errors...")

		exportErrors, err := ths.getErrors(exportTime, "abonents", "abonents_new", domain.SormCustomersRecord{})
		if err != nil {
			return nil, err
		}

		errorsDetected = len(exportErrors)

		var errorCustomers []models.SormCustomersErrorsData
		for _, ec := range exportErrors {
			record, ok := ec.(domain.SormCustomersRecord)
			if ok {
				errorCustomers = append(errorCustomers, ths.toSormCustomersError(record, "unknown reason", customers))
			} else {
				ths.log.Error("Can't cast to record: %+v", ec)
			}
		}

		var sormCustomersToUpsert []models.SormCustomersData

		if len(errorCustomers) > 0 {
			ths.log.Warn("Error customers detected: %d", len(errorCustomers))

			err = ths.sormCustomersErrorsRepo.Create(errorCustomers)
			if err != nil {
				return nil, fmt.Errorf("can't create error customers: %w", err)
			}
		}

		for _, sc := range sormCustomersDiff {
			if !ths.findByLoginInSormCustomersErrors(errorCustomers, sc) {
				sormCustomersToUpsert = append(sormCustomersToUpsert, sc)
			}
		}

		ths.log.Info("sormCustomersToUpsert=%d", len(sormCustomersToUpsert))

		if len(sormCustomersToUpsert) > 0 {
			err = ths.sormCustomersRepo.Upsert(sormCustomersToUpsert)
			if err != nil {
				return nil, fmt.Errorf("can't upsert SORM customers: %w", err)
			}
		}
	}

	var status = "OK"
	if errorsDetected > 0 {
		status = "ERRORS"
	}

	err = ths.sormExportStatusRepo.Upsert(models.SormExportStatusData{
		FileName:    domain.SormCustomersRecord{}.FileName(),
		ExportCount: recordsExported,
		Errors:      errorsDetected,
		Status:      status,
	})
	if err != nil {
		return nil, fmt.Errorf("can't set export status: %w", err)
	}

	return records, nil
}

func (ths *Exporter) findByLoginInSormCustomersErrors(sormCustomers []models.SormCustomersErrorsData, c models.SormCustomersData) bool {
	for _, sc := range sormCustomers {
		if sc.Login == c.Login {
			return true
		}
	}
	return false
}

func (ths *Exporter) findInSormCustomers(sormCustomers []models.SormCustomersData, c models.SormCustomersData) bool {
	for _, sc := range sormCustomers {
		if sc.Hash == c.Hash {
			return true
		}
	}
	return false
}

func (ths *Exporter) toSormCustomersError(r domain.SormCustomersRecord, reason string, customers []models.CustomerData) models.SormCustomersErrorsData {
	for _, c := range customers {
		if c.Login == r.Login {
			return models.SormCustomersErrorsData{
				Login:      r.Login,
				CustomerId: models.NewNullUUID(c.Id),
				Reason:     reason,
			}
		}
	}
	return models.SormCustomersErrorsData{}
}

func (ths *Exporter) fromSormCustomerData(c models.SormCustomersData) domain.SormCustomersRecord {
	return domain.SormCustomersRecord{
		OrgUnit:             "1",
		Login:               c.Login,
		Email:               c.Email,
		Phone:               c.Phone,
		ContractDate:        c.ContractDate,
		ContractNumber:      c.ContractNumber,
		Status:              c.Status,
		Start:               c.Start,
		End:                 c.End,
		Type:                c.Type,
		NameType:            c.NameType,
		Name:                c.Name,
		DocType:             c.DocType,
		Doc:                 c.Doc,
		DocCode:             c.DocCode,
		RegistrationType:    c.RegistrationType,
		RegistrationAddress: c.RegistrationAddress,
		DeviceType:          c.DeviceType,
		DeviceAddress:       c.DeviceAddress,
		PostType:            c.PostType,
		PostAddress:         c.PostAddress,
		BillType:            c.BillType,
		BillAddress:         c.BillAddress,
	}
}

func (ths *Exporter) fromCustomerData(c models.CustomerData) domain.SormCustomersRecord {
	docDate := c.Contract.DocumentDate.Format("02.01.2006")
	address := strings.Join([]string{c.Street.Name, c.Build, c.Flat}, " ")
	doc := fmt.Sprintf("%s %s %s", c.Contract.Document, docDate, c.Contract.DocumentIssuedBy)
	contractDate := c.Contract.Date.Format("02.01.2006")

	var ip string
	if len(c.ServiceInternetCustomData) > 0 {
		ip = c.ServiceInternetCustomData[0].Ip
	}

	record := domain.SormCustomersRecord{
		OrgUnit:             "1",
		IP:                  ip,
		Login:               c.Login,
		Email:               c.Email,
		Phone:               c.Phone,
		ContractDate:        contractDate,
		BirthDate:           c.Contract.BirthDate.Format(DateFormat),
		ContractNumber:      c.Contract.Number,
		Status:              "0", // 0 - подключен, 1 - отключен
		Start:               c.CreateTs.Format(DateFormat),
		End:                 "",
		Type:                "0", // 0 - физик, 1 - юрик
		NameType:            "1", // 0 - структурированные, 1 - неструктурированные
		Name:                c.Name,
		DocType:             "1", // 0 - структурированные, 1 - неструктурированные
		Doc:                 doc,
		DocCode:             "000",
		RegistrationType:    "1", // 0 - структурированные, 1 - неструктурированные
		RegistrationAddress: address,
		DeviceType:          "1", // 0 - структурированные, 1 - неструктурированные
		DeviceAddress:       address,
		PostType:            "1", // 0 - структурированные, 1 - неструктурированные
		PostAddress:         address,
		BillType:            "1", // 0 - структурированные, 1 - неструктурированные
		BillAddress:         address,
	}

	if c.Contract.Type == domain.CustomerContractTypeEnterprise {
		record.EnterpriseName = c.Name
		record.EnterpriseINN = c.Contract.INN
		if len(c.BankAccounts) > 0 {
			record.EnterpriseBank = c.BankAccounts[0].BankName
			record.EnterpriseBankAccount = c.BankAccounts[0].Number
		}
	}

	return record
}
