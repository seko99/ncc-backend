package exporter

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/domain"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
)

func (ths *Exporter) exportGates() ([]exporter.ExportData, error) {
	gates, err := ths.gatewayRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get gates: %w", err)
	}

	var records []exporter.ExportData

	for _, s := range gates {
		records = append(records, domain.SormGatesRecord{
			OrgUnit:  "1",
			IP:       s.IP,
			Start:    s.CreateTs.Format(DateFormat),
			End:      "",
			Descr:    s.Descr,
			Country:  s.Country,
			Region:   s.Region,
			City:     s.City,
			District: s.District,
			Street:   s.Street,
			Build:    s.Build,
			Type:     s.Type,
		})
	}

	err = ths.exportData(records)
	if err != nil {
		return nil, fmt.Errorf("can't export gates: %w", err)
	}

	err = ths.sormExportStatusRepo.Upsert(models.SormExportStatusData{
		FileName:    domain.SormGatesRecord{}.FileName(),
		ExportCount: len(records),
		Errors:      0,
		Status:      "OK",
	})
	if err != nil {
		return nil, fmt.Errorf("can't set export status: %w", err)
	}

	return records, nil
}
