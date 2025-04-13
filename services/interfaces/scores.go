package interfaces

//go:generate mockgen -destination=mocks/mock_scores_service.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/interfaces Scores

type Scores interface {
	Process(dryRun bool) error
}
