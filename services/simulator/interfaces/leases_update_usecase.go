package interfaces

//go:generate mockgen -destination=../mocks/mock_leases_update_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces LeasesUpdateUsecase

type LeasesUpdateUsecase interface {
	Execute() error
}
