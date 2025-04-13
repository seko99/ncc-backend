package interfaces

//go:generate mockgen -destination=../mocks/mock_issues_delete_all_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces IssuesDeleteAllUsecase

type IssuesDeleteAllUsecase interface {
	Execute() error
}
