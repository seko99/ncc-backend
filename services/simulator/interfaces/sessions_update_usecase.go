package interfaces

//go:generate mockgen -destination=../mocks/mock_sessions_update_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces SessionsUpdateUsecase

type SessionsUpdateUsecase interface {
	Execute() error
}
