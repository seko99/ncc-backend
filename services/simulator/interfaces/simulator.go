package interfaces

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/services/simulator"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
)

//go:generate mockgen -destination=../mocks/mock_simulator_service.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces Simulator

type Simulator interface {
	GetSessionCache() ([]models.SessionData, error)
	GetLeasesCache() ([]models.LeaseData, error)
	GetCustomerCache() ([]models.CustomerData, error)
	GetNASCache() ([]models.NasData, error)

	CreateFakeData(req dto.FakeDataCreateUsecaseRequest) error
	ClearFakeData(req dto.FakeDataClearUsecaseRequest) error
	InitDictionaries() error
	Cleanup() error
	UpdateLeases() error
	UpdateMap() error
	UpdateSessions() error
	CreateIssues() error
	DeleteIssues() error
	DropSessions() error

	SetBRASParams(params simulator.BrasParams)

	StartRadiusSessions(req dto.RadiusUsecaseRequest) (dto.RadiusUsecaseResponse, error)
	StopRadiusSessions(req dto.RadiusUsecaseRequest) error
	UpdateRadiusSessions(req dto.RadiusUsecaseRequest) error
	KillRadiusSessions(req dto.RadiusKillSessionsUsecaseRequest) (dto.RadiusKillSessionsUsecaseResponse, error)
}
