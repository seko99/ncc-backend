package interfaces

//go:generate mockgen -destination=mocks/mock_exporter_service.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/interfaces Exporter

type Exporter interface {
	Run() error
}
