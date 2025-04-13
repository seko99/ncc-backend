package domain

import (
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
)

type SormPayTypesRecord struct {
	OrgUnit     string
	ServiceCode string
	Start       string
	End         string
	Descr       string
}

func (s SormPayTypesRecord) FileName() string {
	return "dictionaries/pay-types/pay_types"
}

func (s SormPayTypesRecord) Header() []string {
	return []string{
		"org_unit",
		"code",
		"start",
		"end",
		"descr",
	}
}

func (s SormPayTypesRecord) FromSlice(data []string) (exporter.ExportData, error) {
	var record SormPayTypesRecord

	if len(data) != 5 {
		return record, fmt.Errorf("wrong data: %v", data)
	}

	record.OrgUnit = data[0]
	record.ServiceCode = data[1]
	record.Start = data[2]
	record.End = data[3]
	record.Descr = data[4]

	return record, nil
}

func (s SormPayTypesRecord) ToSlice() []string {
	return []string{
		s.OrgUnit,
		s.ServiceCode,
		s.Start,
		s.End,
		s.Descr,
	}
}
