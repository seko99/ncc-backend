package domain

import (
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
)

type SormPaymentsRecord struct {
	OrgUnit     string
	PaymentCode string
	Login       string
	IP          string
	Date        string
	Amount      string
	Descr       string
}

func (s SormPaymentsRecord) FileName() string {
	return "payments/balance-fillup/balance-fillup"
}

func (s SormPaymentsRecord) Header() []string {
	return []string{
		"org_unit",
		"payment_code",
		"login",
		"ip",
		"date",
		"amount",
		"descr",
	}
}

func (s SormPaymentsRecord) FromSlice(data []string) (exporter.ExportData, error) {
	var record SormPaymentsRecord

	if len(data) != 7 {
		return record, fmt.Errorf("wrong data: %v", data)
	}

	record.OrgUnit = data[0]
	record.PaymentCode = data[1]
	record.Login = data[2]
	record.IP = data[3]
	record.Date = data[4]
	record.Amount = data[5]
	record.Descr = data[6]

	return record, nil
}

func (s SormPaymentsRecord) ToSlice() []string {
	return []string{
		s.OrgUnit,
		s.PaymentCode,
		s.Login,
		s.IP,
		s.Date,
		s.Amount,
		s.Descr,
	}
}
