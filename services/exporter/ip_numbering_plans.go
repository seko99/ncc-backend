package exporter

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/domain"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"fmt"
)

func (ths *Exporter) exportNumberingPlan() ([]exporter.ExportData, error) {
	status, err := ths.sormExportStatusRepo.GetByFileName(domain.SormIpNumberingRecord{}.FileName())
	if err != nil {
		return nil, fmt.Errorf("can't get export status: %w", err)
	}

	if status.Status != "PENDING" {
		return []exporter.ExportData{}, nil
	}

	ipNumbering, err := ths.ipNumberingRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get IP numbering: %w", err)
	}

	var records []exporter.ExportData

	for _, s := range ipNumbering {
		records = append(records, domain.SormIpNumberingRecord{
			OrgUnit: "1",
			Network: s.Network,
			Mask:    s.Mask,
			Start:   s.CreateTs.Format(DateFormat),
			End:     "",
			Descr:   s.Name,
		})
	}

	err = ths.exportData(records)
	if err != nil {
		return nil, fmt.Errorf("can't export IP numbering plan: %w", err)
	}

	err = ths.sormExportStatusRepo.Upsert(models.SormExportStatusData{
		FileName:    domain.SormIpNumberingRecord{}.FileName(),
		ExportCount: len(records),
		Errors:      0,
		Status:      "OK",
	})
	if err != nil {
		return nil, fmt.Errorf("can't set export status: %w", err)
	}

	return records, nil
}
