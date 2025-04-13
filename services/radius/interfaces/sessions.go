package interfaces

import (
	"code.evixo.ru/ncc/ncc-backend/services/radius/dto"
)

//go:generate mockgen -destination=../mocks/mock_session_start_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/radius/interfaces SessionStartUsecase
//go:generate mockgen -destination=../mocks/mock_session_stop_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/radius/interfaces SessionStopUsecase
//go:generate mockgen -destination=../mocks/mock_session_update_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/radius/interfaces SessionUpdateUsecase
//go:generate mockgen -destination=../mocks/mock_session_watcher_usecase.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/radius/interfaces SessionWatcherUsecase

type SessionStartUsecase interface {
	Execute(req dto.SessionStartRequest) (*dto.SessionStartResponse, error)
}

type SessionStopUsecase interface {
	Execute(req dto.SessionStopRequest) (*dto.SessionStopResponse, error)
}

type SessionUpdateUsecase interface {
	Execute(req dto.SessionUpdateRequest) (*dto.SessionUpdateResponse, error)
}

type SessionWatcherUsecase interface {
	Execute() error
}
