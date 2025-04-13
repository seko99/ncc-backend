package exporter

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/domain"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
)

func (ths *Exporter) exportServices() ([]exporter.ExportData, error) {
	services, err := ths.serviceInternetRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get serviceInternet: %w", err)
	}

	var records []exporter.ExportData

	for _, s := range services {
		records = append(records, domain.SormSupplementServicesRecord{
			OrgUnit:     "1",
			ServiceCode: s.Code,
			ServiceName: s.Name,
			Start:       s.CreateTs.Format(DateFormat),
			End:         "",
			Descr:       s.Name,
		})
	}

	err = ths.exportData(records)
	if err != nil {
		return nil, fmt.Errorf("can't export services: %w", err)
	}

	err = ths.sormExportStatusRepo.Upsert(models.SormExportStatusData{
		FileName:    domain.SormSupplementServicesRecord{}.FileName(),
		ExportCount: len(records),
		Errors:      0,
		Status:      "OK",
	})
	if err != nil {
		return nil, fmt.Errorf("can't set export status: %w", err)
	}

	return records, nil
}
