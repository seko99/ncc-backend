package interfaces

import "code.evixo.ru/ncc/ncc-backend/services/simulator/dto"

//go:generate mockgen -destination=../mocks/mock_bras_get_stat_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces BrasGetStatUsecase

type BrasGetStatUsecase interface {
	Execute(req dto.BrasGetStatUsecaseRequest) (dto.BrasGetStatUsecaseResponse, error)
}
