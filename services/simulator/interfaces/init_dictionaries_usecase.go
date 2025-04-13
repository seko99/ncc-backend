package interfaces

//go:generate mockgen -destination=../mocks/mock_init_dictionaries_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces InitDictionariesUsecase

type InitDictionariesUsecase interface {
	Execute() error
}
