package exporter

import "time"

//go:generate mockgen -destination=mocks/mock_exporter_writer.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter ExportWriter
//go:generate mockgen -destination=mocks/mock_exporter_data.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter ExportData

type ExportData interface {
	Header() []string
	ToSlice() []string
	FromSlice(data []string) (ExportData, error)
	FileName() string
}

type ExportWriter interface {
	Write(data []ExportData, withHeader ...bool) error
	GetErrors(exportTime time.Time, path, errorFileName string, d ExportData) ([]ExportData, error)
}
