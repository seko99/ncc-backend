package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/services/radius/dto"
	"code.evixo.ru/ncc/ncc-backend/services/radius/interfaces"
	"fmt"
	"layeh.com/radius/rfc2866"
	"time"
)

const (
	DefaultSessionTimeout = 180
)

type SessionWatcherUsecase struct {
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
	stopUsecase     interfaces.SessionStopUsecase
}

func (ths *SessionWatcherUsecase) Execute() error {
	sessions, err := ths.sessionCache.Get()
	if err != nil {
		return fmt.Errorf("can't get sessions: %w", err)
	}
	for _, s := range sessions {
		var timeout time.Duration

		if s.Nas.SessionTimeout > 0 {
			timeout = time.Duration(s.Nas.SessionTimeout)
		} else {
			timeout = time.Duration(DefaultSessionTimeout)
		}

		if time.Now().After(s.LastAlive.Add(timeout * time.Second)) {
			ts := time.Now()
			_, err := ths.stopUsecase.Execute(dto.SessionStopRequest{
				AcctSessionId:      s.AcctSessionId,
				AcctTerminateCause: uint32(rfc2866.AcctTerminateCause_Value_SessionTimeout),
			})
			if err != nil {
				ths.log.Error("Can't stop session: %v", err)
			} else {
				ths.log.Info("Session stopped by timeout in %v", time.Since(ts))
			}
		}
	}
	return nil
}

func NewSessionWatcherUsecase(
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
	stopUsecase interfaces.SessionStopUsecase,
) SessionWatcherUsecase {
	return SessionWatcherUsecase{
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
		stopUsecase:     stopUsecase,
	}
}
