package interfaces

import "code.evixo.ru/ncc/ncc-backend/services/simulator/dto"

//go:generate mockgen -destination=../mocks/mock_radius_start_all_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces RadiusStartAllUsecase

type RadiusStartAllUsecase interface {
	Execute(req dto.RadiusUsecaseRequest) (dto.RadiusUsecaseResponse, error)
}
