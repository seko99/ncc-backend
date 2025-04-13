package simulator

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"context"
	"fmt"
	"github.com/google/uuid"
	rad "layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"
	"math/rand"
	"net"
	"sync"
	"time"
)

type RadiusServerData struct {
	Secret     string
	IP         string
	Identifier string
	Auth       string
	Acct       string
}

func (ths *Simulator) StartRadiusSessions(req dto.RadiusUsecaseRequest) (dto.RadiusUsecaseResponse, error) {
	response := dto.RadiusUsecaseResponse{}

	leases, err := ths.leasesCache.Get()
	if err != nil {
		return response, fmt.Errorf("can't get leases: %w", err)
	}
	response.Leases = len(leases)

	persistentSessions, err := ths.sessionsRepo.Get()
	if err != nil {
		return response, fmt.Errorf("can't get sessions for cache: %w", err)
	}
	response.PersistentSessions = len(persistentSessions)

	err = ths.sessionCache.DeleteAll()
	if err != nil {
		return response, fmt.Errorf("can't delete all cached sessions: %w", err)
	}

	err = ths.sessionCache.Create(persistentSessions)
	if err != nil {
		return response, fmt.Errorf("can't init session cache: %w", err)
	}

	nases, err := ths.nasesCache.Get()
	if err != nil || len(nases) == 0 {
		return response, fmt.Errorf("can't get NASes: %w", err)
	}
	nas := nases[0]

	limit := len(leases)
	if req.Limit > 0 {
		limit = req.Limit
	}

	var nasPort uint32 = 1000000

	wg := sync.WaitGroup{}
	for _, lease := range leases {
		l := lease

		if limit < len(leases) && l.Customer.ServiceInternetState != models.ServiceStateEnabled {
			continue
		}

		if _, err := ths.sessionCache.GetByIP(l.Ip); err == nil {
			continue
		}

		wg.Add(1)
		go func() {
			packet := rad.New(rad.CodeAccessRequest, []byte(req.Secret))
			rfc2865.UserName_SetString(packet, l.Ip)
			rfc2865.UserPassword_SetString(packet, l.Customer.Password)
			rfc2865.NASIPAddress_Set(packet, net.ParseIP(req.NasIP))
			rfc2865.NASIdentifier_SetString(packet, req.NasIdentifier)
			rfc2865.NASPortType_Set(packet, rfc2865.NASPortType_Value_Ethernet)
			radResponse, err := rad.Exchange(context.Background(), packet, req.Auth)
			if err != nil || radResponse == nil {
				ths.log.Error("Exchange error: %v", err)
				return
			}

			nasPort++
			response.Sent++

			switch radResponse.Code {
			case rad.CodeAccessAccept:
				response.Accepted++

				sessionId := uuid.NewString()

				data := []models.SessionData{
					{
						StartTime:     time.Now(),
						LastAlive:     time.Now(),
						AcctSessionId: sessionId,
						CustomerId:    l.CustomerId,
						Login:         l.Customer.Login,
						Ip:            l.Ip,
						Nas:           nas,
						NasId:         models.NewNullUUID(nas.Id),
						NasPort:       nasPort,
					},
				}
				err := ths.sessionCache.Create(data)
				if err != nil {
					ths.log.Error("Can't create session: %v", err)
				}

				packet := rad.New(rad.CodeAccountingRequest, []byte(req.Secret))
				rfc2866.AcctStatusType_Set(packet, rfc2866.AcctStatusType_Value_Start)
				rfc2865.ServiceType_Add(packet, rfc2865.ServiceType_Value_FramedUser)
				rfc2865.FramedProtocol_Add(packet, rfc2865.FramedProtocol_Value_PPP)
				rfc2865.UserName_SetString(packet, l.Ip)
				rfc2865.NASIPAddress_Set(packet, net.ParseIP(req.NasIP))
				rfc2865.NASIdentifier_AddString(packet, req.NasIdentifier)
				rfc2865.NASPortType_Add(packet, rfc2865.NASPortType_Value_Ethernet)
				rfc2865.NASPort_Add(packet, rfc2865.NASPort(nasPort))
				//rfc2865.CallingStationID_Add(packet, net.ParseIP("192.168.88.239"))

				rfc2866.AcctAuthentic_Add(packet, rfc2866.AcctAuthentic_Value_RADIUS)
				rfc2866.AcctSessionID_Add(packet, []byte(sessionId))
				rfc2866.AcctSessionTime_Set(packet, 0)

				radResponse, err = rad.Exchange(context.Background(), packet, req.Acct)
				if err != nil || radResponse == nil {
					ths.log.Error("Exchange error: %v", err)
				}
			case rad.CodeAccessReject:
				response.Rejected++
			default:
				ths.log.Warn("Unknown dode: %v (%s)", radResponse.Code, l.Ip)
			}
			wg.Done()
		}()

		limit--
		if limit <= 0 {
			break
		}
	}

	wg.Wait()

	cache, err := ths.sessionCache.Get()
	if err != nil {
		return response, fmt.Errorf("can't get cache: %w", err)
	}

	response.CacheSessions = len(cache)

	ths.log.Info("Accepted/rejected/sent/cache/leases: %d/%d/%d/%d/%d", response.Accepted, response.Rejected, response.Sent, response.CacheSessions, len(leases))
	return response, nil
}

func (ths *Simulator) radiusInterimUpdate(req dto.RadiusUsecaseRequest) {
	for {
		if !ths.brasParams.SendInterims {
			continue
		}

		sessions, err := ths.sessionCache.Get()
		if err != nil {
			ths.log.Error("Can't get sessions: %v", err)
		}

		for _, session := range sessions {
			interval := time.Duration(session.Nas.InterimInterval)
			t := time.Now()
			add := session.LastAlive.Add(interval * time.Second)
			if t.After(add) {
				session.LastAlive = time.Now()
				session.Duration = time.Now().Unix() - session.StartTime.Unix()
				session.OctetsIn += int64(rand.Intn(1024 * 100))
				session.OctetsOut += int64(rand.Intn(1024 * 10))
				ths.sendAcctInterimUpdate(RadiusServerData{
					Secret:     req.Secret,
					IP:         req.NasIP,
					Identifier: req.NasIdentifier,
					Auth:       req.Auth,
					Acct:       req.Acct,
				}, session)
				err := ths.sessionCache.UpdateBySessionId(session)
				if err != nil {
					ths.log.Error("Can't update session: %v", err)
				}
			}
		}

		time.Sleep(time.Second)
	}
}

func (ths *Simulator) sendAcctInterimUpdate(srv RadiusServerData, session models.SessionData) {
	packet := rad.New(rad.CodeAccountingRequest, []byte(srv.Secret))
	rfc2866.AcctStatusType_Set(packet, rfc2866.AcctStatusType_Value_InterimUpdate)
	rfc2865.ServiceType_Add(packet, rfc2865.ServiceType_Value_FramedUser)
	rfc2865.FramedIPAddress_Add(packet, net.ParseIP(session.Ip))
	rfc2865.FramedProtocol_Add(packet, rfc2865.FramedProtocol_Value_PPP)
	rfc2865.UserName_SetString(packet, session.Login)
	rfc2865.NASIPAddress_Set(packet, net.ParseIP(srv.IP))
	rfc2865.NASIdentifier_AddString(packet, srv.Identifier)
	rfc2865.NASPortType_Add(packet, rfc2865.NASPortType_Value_Ethernet)
	rfc2865.NASPort_Add(packet, rfc2865.NASPort(session.NasPort))

	rfc2866.AcctInputOctets_Set(packet, rfc2866.AcctInputOctets(session.OctetsIn))
	rfc2866.AcctInputPackets_Set(packet, rfc2866.AcctInputPackets(session.OctetsOut))
	rfc2866.AcctOutputOctets_Set(packet, rfc2866.AcctOutputOctets(session.OctetsIn/1024))
	rfc2866.AcctOutputPackets_Set(packet, rfc2866.AcctOutputPackets(session.OctetsOut/1024))

	//rfc2865.CallingStationID_Add(packet, net.ParseIP("192.168.88.239"))

	rfc2866.AcctAuthentic_Add(packet, rfc2866.AcctAuthentic_Value_RADIUS)
	rfc2866.AcctSessionID_Add(packet, []byte(session.AcctSessionId))
	rfc2866.AcctSessionTime_Set(packet, rfc2866.AcctSessionTime(session.Duration))

	radResponse, err := rad.Exchange(context.Background(), packet, srv.Acct)
	if err != nil || radResponse == nil {
		ths.log.Error("Exchange error: %v", err)
	}
}

func (ths *Simulator) KillRadiusSessions(req dto.RadiusKillSessionsUsecaseRequest) (dto.RadiusKillSessionsUsecaseResponse, error) {
	sessions, err := ths.sessionCache.Get()
	if err != nil {
		return dto.RadiusKillSessionsUsecaseResponse{}, fmt.Errorf("can't get sessions: %w", err)
	}

	if req.Sessions > len(sessions) {
		return dto.RadiusKillSessionsUsecaseResponse{}, fmt.Errorf("sessions should be <= len(sessionCache) [%d vs %d]", req.Sessions, len(sessions))
	}

	killed := 0
	for i := 0; i < req.Sessions; i++ {
		sessionIndex := i
		if req.Random {
			sessionIndex = rand.Intn(len(sessions) - 1)
		}
		err := ths.sessionCache.DeleteBySessionId(sessions[sessionIndex])
		if err != nil {
			ths.log.Error("Can't delete session: %v", err)
		}
		killed++
		sessions, err = ths.sessionCache.Get()
		if err != nil {
			ths.log.Error("Can't get sessions while deleting:: %v", err)
		}
	}

	return dto.RadiusKillSessionsUsecaseResponse{
		Killed: killed,
	}, nil
}

func (ths *Simulator) UpdateRadiusSessions(req dto.RadiusUsecaseRequest) error {
	sessions, err := ths.sessionCache.Get()
	if err != nil {
		return fmt.Errorf("can't get sessions: %w", err)
	}

	for _, session := range sessions {
		session.LastAlive = time.Now()
		session.Duration = time.Now().Unix() - session.StartTime.Unix()
		session.OctetsIn += int64(rand.Intn(1024 * 100))
		session.OctetsOut += int64(rand.Intn(1024 * 10))
		ths.sendAcctInterimUpdate(RadiusServerData{
			Secret:     req.Secret,
			IP:         req.NasIP,
			Identifier: req.NasIdentifier,
			Auth:       req.Auth,
			Acct:       req.Acct,
		}, session)
		err := ths.sessionsRepo.Update(session)
		if err != nil {
			ths.log.Error("Can't update session: %v", err)
		}
	}

	ths.log.Info("Sent %d interims", len(sessions))

	return nil
}

func (ths *Simulator) StopRadiusSessions(req dto.RadiusUsecaseRequest) error {
	var nasPort uint32 = 1000000

	sessions, err := ths.sessionCache.Get()
	if err != nil {
		return fmt.Errorf("can't get sessions: %w", err)
	}

	for _, session := range sessions {
		packet := rad.New(rad.CodeAccountingRequest, []byte(req.Secret))
		rfc2866.AcctStatusType_Set(packet, rfc2866.AcctStatusType_Value_Stop)
		rfc2865.ServiceType_Add(packet, rfc2865.ServiceType_Value_FramedUser)
		rfc2865.FramedProtocol_Add(packet, rfc2865.FramedProtocol_Value_PPP)
		rfc2865.UserName_SetString(packet, session.Login)
		rfc2865.NASIPAddress_Add(packet, net.ParseIP(req.NasIP))
		rfc2865.NASIdentifier_AddString(packet, req.NasIdentifier)
		rfc2865.NASPortType_Add(packet, rfc2865.NASPortType_Value_Ethernet)
		rfc2865.NASPort_Add(packet, rfc2865.NASPort(nasPort))
		//rfc2865.CallingStationID_Add(packet, net.ParseIP("192.168.88.239"))

		rfc2866.AcctAuthentic_Add(packet, rfc2866.AcctAuthentic_Value_RADIUS)
		rfc2866.AcctSessionID_Add(packet, []byte(session.AcctSessionId))
		rfc2866.AcctSessionTime_Set(packet, rfc2866.AcctSessionTime(time.Now().Unix()-session.StartTime.Unix()))
		rfc2866.AcctTerminateCause_Set(packet, rfc2866.AcctTerminateCause_Value_UserRequest)

		response, err := rad.Exchange(context.Background(), packet, req.Acct)
		if err != nil || response == nil {
			ths.log.Error("Exchange error: %v", err)
		}
	}

	err = ths.sessionCache.DeleteAll()
	if err != nil {
		return fmt.Errorf("can't delete sessions: %w", err)
	}

	return nil
}
