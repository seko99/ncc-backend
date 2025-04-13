package interfaces

//go:generate mockgen -destination=../mocks/mock_issues_create_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces IssuesCreateUsecase

type IssuesCreateUsecase interface {
	Execute() error
}
