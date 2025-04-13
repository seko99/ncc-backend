package interfaces

import "code.evixo.ru/ncc/ncc-backend/services/simulator/dto"

//go:generate mockgen -destination=../mocks/mock_bras_get_sessions_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces BrasGetSessionsUsecase

type BrasGetSessionsUsecase interface {
	Execute(req dto.BrasGetSessionsUsecaseRequest) (dto.BrasGetSessionsUsecaseResponse, error)
}
