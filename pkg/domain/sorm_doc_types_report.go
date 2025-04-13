package domain

import (
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
)

type SormDocTypesRecord struct {
	OrgUnit string
	Code    string
	Start   string
	End     string
	Descr   string
}

func (s SormDocTypesRecord) FileName() string {
	return "dictionaries/doc-types/doc_types"
}

func (s SormDocTypesRecord) Header() []string {
	return []string{
		"org_unit",
		"code",
		"start",
		"end",
		"descr",
	}
}

func (s SormDocTypesRecord) FromSlice(data []string) (exporter.ExportData, error) {
	var record SormDocTypesRecord

	if len(data) != 5 {
		return record, fmt.Errorf("wrong data: %v", data)
	}

	record.OrgUnit = data[0]
	record.Code = data[1]
	record.Start = data[2]
	record.End = data[3]
	record.Descr = data[4]

	return record, nil
}

func (s SormDocTypesRecord) ToSlice() []string {
	return []string{
		s.OrgUnit,
		s.Code,
		s.Start,
		s.End,
		s.Descr,
	}
}
