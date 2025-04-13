package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/services/radius/dto"
	"fmt"
	"time"
)

type SessionStopUsecase struct {
	cfg             *config.Config
	log             logger.Logger
	sessionsRepo    repository.Sessions
	sessionsLogRepo repository.SessionsLog
	leasesRepo      repository.DhcpLeases
	customersRepo   repository.Customers
	nasesRepo       repository.Nases
	sessionCache    repository.Sessions
	leasesCache     repository.DhcpLeases
	customersCache  repository.Customers
	nasesCache      repository.Nases
}

func (ths *SessionStopUsecase) Execute(req dto.SessionStopRequest) (*dto.SessionStopResponse, error) {
	sessionId := req.AcctSessionId

	reqSession := models.SessionData{
		AcctSessionId: sessionId,
		LastAlive:     time.Now(),
		Duration:      req.AcctSessionTime,
		OctetsIn:      int64(req.AcctInputOctets),
		OctetsOut:     int64(req.AcctOutputOctets),
	}

	session, err := ths.sessionCache.GetBySessionId(reqSession.AcctSessionId)
	if err != nil {
		return nil, fmt.Errorf("can't get session: %w", err)
	}

	if req.AcctInputOctets != uint32(session.OctetsIn) ||
		req.AcctOutputOctets != uint32(session.OctetsOut) ||
		req.AcctSessionTime != session.Duration {

		err := ths.sessionsRepo.UpdateBySessionId(reqSession)
		if err != nil {
			return nil, fmt.Errorf("can't update session: %w", err)
		}
	}

	err = ths.sessionsRepo.DeleteBySessionId(models.SessionData{
		AcctSessionId: sessionId,
	})
	if err != nil {
		return nil, fmt.Errorf("can't delete session: %w", err)
	}

	sessionLog := models.SessionsLogData{}.FromSession(session)
	sessionLog.TerminateCause = models.TerminateCause(req.AcctTerminateCause).String()
	err = ths.sessionsLogRepo.Create(sessionLog)
	if err != nil {
		return nil, fmt.Errorf("can't create session log: %w", err)
	}

	return &dto.SessionStopResponse{
		AcctSessionId:    sessionId,
		AcctSessionTime:  reqSession.Duration,
		AcctInputOctets:  uint32(reqSession.OctetsIn),
		AcctOutputOctets: uint32(reqSession.OctetsOut),
	}, nil
}

func NewSessionStopUsecase(
	cfg *config.Config,
	log logger.Logger,
	sessionsRepo repository.Sessions,
	sesssionsLogRepo repository.SessionsLog,
	leasesRepo repository.DhcpLeases,
	customersRepo repository.Customers,
	nasesRepo repository.Nases,
	sessionCache repository.Sessions,
	leasesCache repository.DhcpLeases,
	customersCache repository.Customers,
	nasesCache repository.Nases,
) SessionStopUsecase {
	return SessionStopUsecase{
		cfg:             cfg,
		log:             log,
		sessionsRepo:    sessionsRepo,
		sessionsLogRepo: sesssionsLogRepo,
		leasesRepo:      leasesRepo,
		customersRepo:   customersRepo,
		nasesRepo:       nasesRepo,
		sessionCache:    sessionCache,
		leasesCache:     leasesCache,
		customersCache:  customersCache,
		nasesCache:      nasesCache,
	}
}
