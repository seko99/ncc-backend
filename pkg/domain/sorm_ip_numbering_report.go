package domain

import (
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
)

type SormIpNumberingRecord struct {
	OrgUnit string
	Descr   string
	Network string
	Mask    string
	Start   string
	End     string
}

func (s SormIpNumberingRecord) FileName() string {
	return "dictionaries/ip-numbering-plan/ip-numbering-plan"
}

func (s SormIpNumberingRecord) Header() []string {
	return []string{
		"org_unit",
		"descr",
		"network",
		"mask",
		"start",
		"end",
	}
}

func (s SormIpNumberingRecord) FromSlice(data []string) (exporter.ExportData, error) {
	var record SormIpNumberingRecord

	if len(data) != 6 {
		return record, fmt.Errorf("wrong data: %v", data)
	}

	record.OrgUnit = data[0]
	record.Descr = data[1]
	record.Network = data[2]
	record.Mask = data[3]
	record.Start = data[4]
	record.End = data[5]

	return record, nil
}

func (s SormIpNumberingRecord) ToSlice() []string {
	return []string{
		s.OrgUnit,
		s.Descr,
		s.Network,
		s.Mask,
		s.Start,
		s.End,
	}
}
