package interfaces

import "code.evixo.ru/ncc/ncc-backend/services/simulator/dto"

//go:generate mockgen -destination=../mocks/mock_fake_data_clear_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces FakeDataClearUsecase

type FakeDataClearUsecase interface {
	Execute(request dto.FakeDataClearUsecaseRequest) error
}
