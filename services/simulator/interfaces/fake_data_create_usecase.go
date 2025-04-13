package interfaces

import "code.evixo.ru/ncc/ncc-backend/services/simulator/dto"

//go:generate mockgen -destination=../mocks/mock_fake_data_create_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces FakeDataCreateUsecase

type FakeDataCreateUsecase interface {
	Execute(request dto.FakeDataCreateUsecaseRequest) error
}
