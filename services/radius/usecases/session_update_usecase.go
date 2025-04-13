package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository/memory"
	"code.evixo.ru/ncc/ncc-backend/services/radius/dto"
	"code.evixo.ru/ncc/ncc-backend/services/radius/interfaces"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type SessionUpdateUsecase struct {
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
	startUsecase    interfaces.SessionStartUsecase
}

func (ths *SessionUpdateUsecase) Execute(req dto.SessionUpdateRequest) (*dto.SessionUpdateResponse, error) {
	sessionId := req.AcctSessionId

	reqSession := models.SessionData{
		AcctSessionId: sessionId,
		Duration:      req.AcctSessionTime,
		LastAlive:     time.Now(),
		OctetsIn:      int64(req.AcctInputOctets),
		OctetsOut:     int64(req.AcctOutputOctets),
		Ip:            req.FramedIp,
	}

	session, err := ths.sessionCache.GetBySessionId(sessionId)
	if err != nil {
		if errors.Is(err, memory.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
			sessionLog, err := ths.sessionsLogRepo.GetBySessionId(sessionId)
			if err != nil {
				if errors.Is(err, memory.ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound) {
					response, err := ths.startSession(req)
					if err != nil {
						return nil, fmt.Errorf("can't start session: %w", err)
					}
					updateResponse := dto.SessionUpdateResponse{}.FromStartResponse(response)
					return &updateResponse, nil
				}
				return nil, fmt.Errorf("can't get session from log: %w", err)
			}
			session, err = ths.moveFromSessionLog(sessionLog)
			if err != nil {
				return nil, fmt.Errorf("can't move from session log: %w", err)
			}
			updateResponse := dto.SessionUpdateResponse{}.FromSession(&session)
			return &updateResponse, nil
		}
		return nil, fmt.Errorf("session find error: %w", err)
	}

	if session.Ip != req.FramedIp {
		return nil, fmt.Errorf("update IP %s not equal to existing session IP %s", req.FramedIp, session.Ip)
	}

	session.Duration = reqSession.Duration
	session.OctetsIn = reqSession.OctetsIn
	session.OctetsOut = reqSession.OctetsOut
	session.LastAlive = reqSession.LastAlive

	err = ths.sessionsRepo.UpdateBySessionId(session)
	if err != nil {
		return nil, fmt.Errorf("can't update session: %w", err)
	}

	return &dto.SessionUpdateResponse{
		AcctSessionId:    sessionId,
		AcctSessionTime:  reqSession.Duration,
		AcctInputOctets:  uint32(reqSession.OctetsIn),
		AcctOutputOctets: uint32(reqSession.OctetsOut),
	}, nil
}

func (ths *SessionUpdateUsecase) startSession(req dto.SessionUpdateRequest) (*dto.SessionStartResponse, error) {

	response, err := ths.startUsecase.Execute(dto.SessionStartRequest{}.FromUpdateRequest(req))
	if err != nil {
		return nil, fmt.Errorf("can't start session: %w", err)
	}

	return response, nil
}

func (ths *SessionUpdateUsecase) moveFromSessionLog(sessionLog models.SessionsLogData) (models.SessionData, error) {
	err := ths.sessionsRepo.Create([]models.SessionData{models.SessionData{}.FromSessionLog(sessionLog)})
	if err != nil {
		return models.SessionData{}, fmt.Errorf("can't create session: %w", err)
	}

	err = ths.sessionsLogRepo.DeleteById(sessionLog.Id)
	if err != nil {
		return models.SessionData{}, fmt.Errorf("can't delete from session log: %w", err)
	}

	return models.SessionData{}, nil
}

func NewSessionUpdateUsecase(
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
	startUsecase interfaces.SessionStartUsecase,
) SessionUpdateUsecase {
	return SessionUpdateUsecase{
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
		startUsecase:    startUsecase,
	}
}
