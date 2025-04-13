package exporter

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/domain"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
)

func (ths *Exporter) exportDocTypes() ([]exporter.ExportData, error) {
	types, err := ths.documentTypesRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get document types: %w", err)
	}

	var records []exporter.ExportData

	for _, s := range types {
		records = append(records, domain.SormDocTypesRecord{
			OrgUnit: "1",
			Code:    s.Code,
			Start:   s.CreateTs.Format(DateFormat),
			End:     "",
			Descr:   s.Name,
		})
	}

	err = ths.exportData(records)
	if err != nil {
		return nil, fmt.Errorf("can't export document types: %w", err)
	}

	err = ths.sormExportStatusRepo.Upsert(models.SormExportStatusData{
		FileName:    domain.SormDocTypesRecord{}.FileName(),
		ExportCount: len(records),
		Errors:      0,
		Status:      "OK",
	})
	if err != nil {
		return nil, fmt.Errorf("can't set export status: %w", err)
	}

	return records, nil
}
