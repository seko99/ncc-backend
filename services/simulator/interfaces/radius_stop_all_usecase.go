package interfaces

import "code.evixo.ru/ncc/ncc-backend/services/simulator/dto"

//go:generate mockgen -destination=../mocks/mock_radius_stop_all_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces RadiusStopAllUsecase

type RadiusStopAllUsecase interface {
	Execute(req dto.RadiusUsecaseRequest) error
}
