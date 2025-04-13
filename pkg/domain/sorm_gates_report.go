package domain

import (
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
)

type SormGatesRecord struct {
	OrgUnit  string
	IP       string
	Start    string
	End      string
	Descr    string
	Country  string
	Region   string
	District string
	City     string
	Street   string
	Build    string
	Type     string
}

func (s SormGatesRecord) FileName() string {
	return "dictionaries/gates/gates"
}

func (s SormGatesRecord) Header() []string {
	return []string{
		"org_unit",
		"ip",
		"start",
		"end",
		"descr",
		"country",
		"region",
		"district",
		"city",
		"street",
		"build",
		"type",
	}
}

func (s SormGatesRecord) FromSlice(data []string) (exporter.ExportData, error) {
	var record SormGatesRecord

	if len(data) != 12 {
		return record, fmt.Errorf("wrong data: %v", data)
	}

	record.OrgUnit = data[0]
	record.IP = data[1]
	record.Start = data[2]
	record.End = data[3]
	record.Descr = data[4]
	record.Country = data[5]
	record.Region = data[6]
	record.District = data[7]
	record.City = data[8]
	record.Street = data[9]
	record.Build = data[10]
	record.Type = data[11]

	return record, nil
}

func (s SormGatesRecord) ToSlice() []string {
	return []string{
		s.OrgUnit,
		s.IP,
		s.Start,
		s.End,
		s.Descr,
		s.Country,
		s.Region,
		s.District,
		s.City,
		s.Street,
		s.Build,
		s.Type,
	}
}
