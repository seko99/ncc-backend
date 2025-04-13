package interfaces

import "code.evixo.ru/ncc/ncc-backend/services/simulator/dto"

//go:generate mockgen -destination=../mocks/mock_radius_update_all_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces RadiusUpdateAllUsecase

type RadiusUpdateAllUsecase interface {
	Execute(req dto.RadiusUsecaseRequest) error
}
