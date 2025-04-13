package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository/memory"
	"code.evixo.ru/ncc/ncc-backend/services/radius/dto"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net"
	"time"
)

const (
	MaxSessions = 1
)

var (
	ErrNotFound    = errors.New("not found")
	ErrDuplicate   = errors.New("duplicate session")
	ErrNoLease     = errors.New("no lease found for IP")
	ErrIPNotEqual  = errors.New("UserName must be equal to IP")
	ErrMaxSessions = errors.New("max sessions reached")
	ErrNoBinding   = errors.New("no binding")
)

type SessionStartUsecase struct {
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

func (ths *SessionStartUsecase) Execute(req dto.SessionStartRequest) (*dto.SessionStartResponse, error) {
	var customer models2.CustomerData

	sessionId := req.AcctSessionId

	if ths.isIP(req.UserName) {
		if req.UserName != req.FramedIp {
			return nil, ErrIPNotEqual
		}

		ip := req.FramedIp

		session, err := ths.sessionCache.GetByIP(ip)
		if err != nil {
			lease, err := ths.leasesCache.GetByIP(ip)
			if err != nil {
				return nil, ErrNoLease
			}

			if !lease.CustomerId.Valid {
				return nil, ErrNoBinding
			}

			customer = lease.Customer
		} else {
			if session.AcctSessionId == sessionId {
				return nil, ErrDuplicate
			}
			err := ths.closePrevSession(session)
			if err != nil {
				return nil, fmt.Errorf("can't close previous session: %w", err)
			}
			customer = session.Customer
		}
	} else {
		sessions, err := ths.sessionCache.GetByLogin(req.UserName)
		if err != nil {
			if !errors.Is(err, memory.ErrNotFound) && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, fmt.Errorf("session find error: %w", err)
			}

			c, err := ths.customersCache.GetByLogin(req.UserName)
			if err != nil {
				return nil, fmt.Errorf("can't find customer: %w", err)
			}

			customer = *c
		} else {
			for _, s := range sessions {
				if s.AcctSessionId == sessionId {
					return nil, ErrDuplicate
				}
			}

			//todo: get MaxSessions from nas or customer settings
			if len(sessions) >= MaxSessions {
				return nil, ErrMaxSessions
			}

			customer = sessions[0].Customer
		}
	}

	if !customer.ServiceInternetId.Valid {
		return nil, fmt.Errorf("no serviceId")
	}

	nas, err := ths.nasesCache.GetByIP(req.NasIpAddress)
	if err != nil {
		return nil, fmt.Errorf("can't find NAS: %w", err)
	}

	newSession := models2.SessionData{
		AcctSessionId:     sessionId,
		Login:             req.UserName,
		StartTime:         time.Now(),
		Duration:          req.AcctSessionTime,
		Ip:                req.FramedIp,
		LastAlive:         time.Now(),
		CustomerId:        models2.NewNullUUID(customer.Id),
		Customer:          customer,
		ServiceInternetId: customer.ServiceInternetId,
		NasId:             models2.NewNullUUID(nas.Id),
		NasName:           nas.Name,
	}

	err = ths.sessionsRepo.Create([]models2.SessionData{newSession})
	if err != nil {
		return nil, fmt.Errorf("can't create session: %w", err)
	}

	return &dto.SessionStartResponse{
		AcctSessionId: sessionId,
	}, nil
}

func (ths *SessionStartUsecase) closePrevSession(session models2.SessionData) error {
	err := ths.sessionsRepo.Delete(session.Id)
	if err != nil {
		return fmt.Errorf("can't delete session: %w", err)
	}

	err = ths.sessionsLogRepo.Create(models2.SessionsLogData{}.FromSession(session))
	if err != nil {
		ths.log.Error("Can't create session log record: %w", err)
	}
	return nil
}

func (ths *SessionStartUsecase) isIP(userName string) bool {
	return net.ParseIP(userName) != nil
}

func NewSessionStartUsecase(
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
) SessionStartUsecase {
	return SessionStartUsecase{
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
