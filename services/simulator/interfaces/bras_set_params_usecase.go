package interfaces

import "code.evixo.ru/ncc/ncc-backend/services/simulator/dto"

//go:generate mockgen -destination=../mocks/mock_bras_set_param_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces BrasSetParamsUsecase

type BrasSetParamsUsecase interface {
	Execute(req dto.BrasSetParamsUsecaseRequest) (dto.BrasSetParamsUsecaseResponse, error)
}
