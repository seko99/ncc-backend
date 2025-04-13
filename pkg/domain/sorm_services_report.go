package domain

import (
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
)

type SormSupplementServicesRecord struct {
	OrgUnit     string
	ServiceCode string
	ServiceName string
	Start       string
	End         string
	Descr       string
}

func (s SormSupplementServicesRecord) FileName() string {
	return "dictionaries/supplement-services/supplement-services"
}

func (s SormSupplementServicesRecord) Header() []string {
	return []string{
		"org_unit",
		"code",
		"name",
		"start",
		"end",
		"descr",
	}
}

func (s SormSupplementServicesRecord) FromSlice(data []string) (exporter.ExportData, error) {
	var record SormSupplementServicesRecord

	if len(data) != 6 {
		return record, fmt.Errorf("wrong data: %v", data)
	}

	record.OrgUnit = data[0]
	record.ServiceCode = data[1]
	record.ServiceName = data[2]
	record.Start = data[3]
	record.End = data[4]
	record.Descr = data[5]
	return record, nil
}

func (s SormSupplementServicesRecord) ToSlice() []string {
	return []string{
		s.OrgUnit,
		s.ServiceCode,
		s.ServiceName,
		s.Start,
		s.End,
		s.Descr,
	}
}
