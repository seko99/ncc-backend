package interfaces

import "code.evixo.ru/ncc/ncc-backend/services/simulator/dto"

//go:generate mockgen -destination=../mocks/mock_radius_kill_sessions_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces RadiusKillSessionsUsecase

type RadiusKillSessionsUsecase interface {
	Execute(req dto.RadiusKillSessionsUsecaseRequest) (dto.RadiusKillSessionsUsecaseResponse, error)
}
