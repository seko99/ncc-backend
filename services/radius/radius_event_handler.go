package radius

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/helpers"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	repository2 "code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository/memory"
	"code.evixo.ru/ncc/ncc-backend/services/radius/dto"
	"code.evixo.ru/ncc/ncc-backend/services/radius/interfaces"
	"code.evixo.ru/ncc/ncc-backend/services/radius/usecases"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gurukami/typ"
	"sync"
	"time"
)

const (
	SessionStart  = "SessionStart"
	SessionStop   = "SessionStop"
	InterimUpdate = "InterimUpdate"

	ConfigRequest  = "ConfigRequest"
	ConfigResponse = "ConfigResponse"

	SessionsRequest  = "SessionsRequest"
	SessionsResponse = "SessionsResponse"

	ReloadRequest  = "ReloadRequest"
	ReloadResponse = "ReloadResponse"

	AccessRequestType  = "AccessRequest"
	AccessResponseType = "AccessResponse"
	AccessAcceptType   = true
	AccessRejectType   = false

	Queue = "radius.events"
)

type User struct {
	Login    string
	Password string
	Ip       []string
	Status   bool
	Attrs    map[string]string
}

type NasAttr struct {
	Attr   string
	Val    string
	Code   uint8
	Vendor uint32
}

type Nas struct {
	Name   string
	Ip     string
	Secret string
	Attrs  []NasAttr
}

type ConfigRequestEvent struct {
}

type ConfigResponseEvent struct {
	Logins []models2.CustomerData `json:"logins"`
}

type SessionMap struct {
	sync.RWMutex
	m map[string]models2.SessionData
}

type NasMap struct {
	sync.RWMutex
	m map[string]models2.NasData
}

type EventHandler struct {
	cfg             *config.Config
	log             logger.Logger
	events          *events.Events
	sessions        *SessionMap
	customersRepo   repository2.Customers
	nasesRepo       repository2.Nases
	sessionsRepo    repository2.Sessions
	sessionsLogRepo repository2.SessionsLog
	leasesRepo      repository2.DhcpLeases
	serviceInternet repository2.ServiceInternet

	sessionCache repository2.Sessions
	leasesCache  repository2.DhcpLeases

	sessionStartUsecase   interfaces.SessionStartUsecase
	sessionStopUsecase    interfaces.SessionStopUsecase
	sessionUpdateUsecase  interfaces.SessionUpdateUsecase
	sessionWatcherUsecase interfaces.SessionWatcherUsecase

	accessMap *AccessMap
	nasMap    *NasMap

	metrics struct {
		sync.RWMutex
		startsQueued   uint
		stopsQueued    uint
		interimsQueued uint
	}

	Wg      sync.WaitGroup
	startup bool
}

func (s *EventHandler) GetSessionCache() repository2.Sessions {
	return s.sessionCache
}

func (s *EventHandler) GetLeasesCache() repository2.DhcpLeases {
	return s.leasesCache
}

func (s *EventHandler) Start() {
	go s.sessionWatcher()
	go s.accessMapUpdater()

	go func() {
		for range time.Tick(time.Second) {
			s.metrics.RLock()
			s.log.Trace("startsQueued=%d, interimQueued=%d", s.metrics.startsQueued, s.metrics.interimsQueued)
			s.metrics.RUnlock()
		}
	}()

	select {}
}

func (s *EventHandler) sessionWatcher() {
	for range time.Tick(time.Second) {
		err := s.sessionWatcherUsecase.Execute()
		if err != nil {
			s.log.Error("Watcher error: %v", err)
		}
	}
}

func (s *EventHandler) getSessionStartRequest(event events.Event) (dto.SessionStartRequest, error) {
	var e dto.SessionStartRequest
	b, _ := json.Marshal(typ.Of(event.Payload).Get("session").Interface().V())
	err := json.Unmarshal(b, &e)
	if err != nil {
		s.log.Error("Can't unmarshal session event: %v", err)
		return dto.SessionStartRequest{}, err
	}

	return e, nil
}

func (s *EventHandler) getSessionStopRequest(event events.Event) (dto.SessionStopRequest, error) {
	var e dto.SessionStopRequest
	b, _ := json.Marshal(typ.Of(event.Payload).Get("session").Interface().V())
	err := json.Unmarshal(b, &e)
	if err != nil {
		s.log.Error("Can't unmarshal session event: %v", err)
		return dto.SessionStopRequest{}, err
	}

	return e, nil
}

func (s *EventHandler) getSessionUpdateRequest(event events.Event) (dto.SessionUpdateRequest, error) {
	var e dto.SessionUpdateRequest
	b, _ := json.Marshal(typ.Of(event.Payload).Get("session").Interface().V())
	err := json.Unmarshal(b, &e)
	if err != nil {
		s.log.Error("Can't unmarshal session event: %v", err)
		return dto.SessionUpdateRequest{}, err
	}

	return e, nil
}

func (s *EventHandler) getAccessItem(login string) (*AccessItem, error) {
	s.accessMap.RLock()
	accessItem, ok := s.accessMap.m[login]
	s.accessMap.RUnlock()

	if !ok {
		return nil, fmt.Errorf("access not found for %s", login)
	}

	return &accessItem, nil
}

func (s *EventHandler) updateAccessMap() {
	customers, err := s.customersRepo.GetByState(models2.CustomerStateActive)
	if err != nil {
		s.log.Error("Can't get customers: %v", err)
	}
	s.log.Debug("Got customers (%d)", len(customers))

	leases, err := s.leasesRepo.Get()
	if err != nil {
		s.log.Error("Can't get leases: %v", err)
	}
	s.log.Debug("Got leases (%d)", len(leases))

	leaseMap := make(map[string]models2.LeaseData)
	for _, l := range leases {
		leaseMap[l.Customer.Login] = l
	}

	accessMap := make(map[string]AccessItem)
	for _, c := range customers {
		lm, hasLease := leaseMap[c.Login]
		var state int
		if c.BlockingState == models2.CustomerStateActive && c.ServiceInternetState == models2.CustomerStateActive {
			state = models2.CustomerStateActive
		}
		accessItem := AccessItem{
			State: state,
			Limits: Limits{
				Login:    c.Login,
				Password: c.Password,
			},
			CustomerId: c.Id,
			ServiceId:  c.ServiceInternetId.UUID.String(),
			HasLease:   hasLease,
		}
		accessMap[c.Login] = accessItem
		accessMap[lm.Ip] = accessItem
	}

	s.accessMap.Lock()
	s.accessMap.m = map[string]AccessItem{}
	for k, v := range accessMap {
		s.accessMap.m[k] = v
	}

	if s.startup {
		if len(s.accessMap.m) > 0 {
			s.startup = false
			s.Wg.Done()
		}
	}

	s.accessMap.Unlock()

	nases, err := s.nasesRepo.Get()
	if err != nil {
		s.log.Error("Can't get NASes: %v", err)
	}

	s.nasMap.Lock()
	s.nasMap.m = map[string]models2.NasData{}
	for _, n := range nases {
		s.nasMap.m[n.Ip] = n
	}
	s.nasMap.Unlock()
}

func (s *EventHandler) accessMapUpdater() {
	for {
		s.updateAccessMap()
		time.Sleep(s.cfg.Radius.Update)
	}
}

func (s *EventHandler) configRequestHandler(event events.Event) map[string]interface{} {
	s.log.Debug("ConfigRequest from %s", event.Publisher.Id)

	nases, err := s.nasesRepo.Get()
	if err != nil {
		s.log.Error("Can't get NASes: %v", err)
	}

	nasMap := map[string]models2.NasData{}
	for _, n := range nases {
		nasMap[n.Ip] = n
	}

	accessMap := map[string]AccessItem{}
	s.accessMap.RLock()
	for k, v := range s.accessMap.m {
		accessMap[k] = v
	}
	s.accessMap.RUnlock()

	return map[string]interface{}{
		"access": accessMap,
		"nases":  nasMap,
	}
}

func (s *EventHandler) interimUpdateHandler(event events.Event) {
	s.log.Info("Updated session in %v", helpers.Timing(func() {
		e, err := s.getSessionUpdateRequest(event)
		if err != nil {
			s.log.Error("Can't get session: %v", err)
			return
		}

		_, err = s.sessionUpdateUsecase.Execute(e)
		if err != nil {
			s.log.Error("can't update session: %v", err)
			return
		}
	}))
}

func (s *EventHandler) sessionStartHandler(event events.Event) {

	s.log.Info("Started session in %v", helpers.Timing(func() {
		e, err := s.getSessionStartRequest(event)
		if err != nil {
			s.log.Error("Can't get session: %v", err)
			return
		}

		_, err = s.sessionStartUsecase.Execute(e)
		if err != nil {
			s.log.Error("can't start session: %v", err)
			return
		}
	}))
}

func (s *EventHandler) sessionStopHandler(event events.Event) {
	s.log.Info("Stopped session in %v", helpers.Timing(func() {
		e, err := s.getSessionStopRequest(event)
		if err != nil {
			s.log.Error("Can't get session: %v", err)
			return
		}

		_, err = s.sessionStopUsecase.Execute(e)
		if err != nil {
			s.log.Error("can't stop session: %v", err)
			return
		}
	}))
}

func NewRadiusEventHandler(
	cfg *config.Config,
	log logger.Logger,
	customersRepo repository2.Customers,
	nasesRepo repository2.Nases,
	leasesRepo repository2.DhcpLeases,
	sessionsRepo repository2.Sessions,
	sessionsLogRepo repository2.SessionsLog,
	serviceInternetRepo repository2.ServiceInternet,
) (*EventHandler, error) {
	radiusEvents, err := events.NewEvents(cfg, log, uuid.NewString(), Queue)
	if err != nil {
		return nil, fmt.Errorf("can't create event system: %v", err)
	}
	radiusEvents.Run()

	sessions := SessionMap{
		m: map[string]models2.SessionData{},
	}

	acccessMap := AccessMap{}
	nasMap := NasMap{}

	r := &EventHandler{
		cfg:             cfg,
		log:             log,
		events:          radiusEvents,
		sessions:        &sessions,
		accessMap:       &acccessMap,
		nasMap:          &nasMap,
		customersRepo:   customersRepo,
		nasesRepo:       nasesRepo,
		sessionsRepo:    sessionsRepo,
		sessionsLogRepo: sessionsLogRepo,
		leasesRepo:      leasesRepo,
		serviceInternet: serviceInternetRepo,
		Wg:              sync.WaitGroup{},
	}

	broadcastEvents, err := events.NewEvents(cfg, log, uuid.NewString(), events.BroadcastEvents)
	if err != nil {
		return nil, fmt.Errorf("can't init event system: %w", err)
	}

	r.sessionCache, err = memory.NewSessions(log, sessionsRepo, broadcastEvents)
	if err != nil {
		return nil, fmt.Errorf("can't init session cache: %w", err)
	}

	customerCache, err := memory.NewCustomers(log, customersRepo, broadcastEvents)
	if err != nil {
		return nil, fmt.Errorf("can't init customer cache: %w", err)
	}

	r.leasesCache, err = memory.NewDhcpLeases(log, leasesRepo, broadcastEvents)
	if err != nil {
		return nil, fmt.Errorf("can't init lease cache: %w", err)
	}

	nasCache, err := memory.NewNases(log, nasesRepo, broadcastEvents)
	if err != nil {
		return nil, fmt.Errorf("can't init NAS cache: %w", err)
	}

	startUsecase := usecases.NewSessionStartUsecase(
		cfg,
		log,
		sessionsRepo,
		sessionsLogRepo,
		leasesRepo,
		customersRepo,
		nasesRepo,
		r.sessionCache,
		r.leasesCache,
		customerCache,
		nasCache,
	)

	sessionStartUsecase := usecases.NewSessionStartUsecase(cfg, log, sessionsRepo, sessionsLogRepo, leasesRepo, customersRepo, nasesRepo, r.sessionCache, r.leasesCache, customerCache, nasCache)
	r.sessionStartUsecase = &sessionStartUsecase

	sessionStopUsecase := usecases.NewSessionStopUsecase(cfg, log, sessionsRepo, sessionsLogRepo, leasesRepo, customersRepo, nasesRepo, r.sessionCache, r.leasesCache, customerCache, nasCache)
	r.sessionStopUsecase = &sessionStopUsecase

	sessionUpdateUsecase := usecases.NewSessionUpdateUsecase(cfg, log, sessionsRepo, sessionsLogRepo, leasesRepo, customersRepo, nasesRepo, r.sessionCache, r.leasesCache, customerCache, nasCache, &startUsecase)
	r.sessionUpdateUsecase = &sessionUpdateUsecase

	sessionWatcherUsecase := usecases.NewSessionWatcherUsecase(cfg, log, sessionsRepo, sessionsLogRepo, leasesRepo, customersRepo, nasesRepo, r.sessionCache, r.leasesCache, customerCache, nasCache, &sessionStopUsecase)
	r.sessionWatcherUsecase = &sessionWatcherUsecase

	err = radiusEvents.SubscribeOnRequest(ConfigRequest, ConfigResponse, r.configRequestHandler)
	if err != nil {
		return nil, fmt.Errorf("can't subscribe on ConfigRequest: %v", err)
	}

	err = radiusEvents.SubscribeOnEvent(InterimUpdate, r.interimUpdateHandler)
	if err != nil {
		return nil, fmt.Errorf("can't subscribe on InterimUpdate: %v", err)
	}

	err = radiusEvents.SubscribeOnEvent(SessionStart, r.sessionStartHandler)
	if err != nil {
		return nil, fmt.Errorf("can't subscribe on SessionStart: %v", err)
	}

	err = radiusEvents.SubscribeOnEvent(SessionStop, r.sessionStopHandler)
	if err != nil {
		return nil, fmt.Errorf("can't subscribe on SessionStop: %v", err)
	}

	r.Wg.Add(1)
	r.startup = true

	return r, nil
}
