package dhcp

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	repository2 "code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository/memory"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

const (
	ConfigRequest  = "ConfigRequest"
	ConfigResponse = "ConfigResponse"

	Queue = "dhcp.events"
)

type EventHandler struct {
	cfg          *config.Config
	log          logger.Logger
	events       *events.Events
	leasesCache  repository2.DhcpLeases
	nasCache     repository2.Nases
	poolsCache   repository2.DhcpPools
	bindingCache repository2.DhcpBindings

	Wg      sync.WaitGroup
	startup bool
}

func (ths *EventHandler) configRequestHandler(event events.Event) map[string]interface{} {
	ths.log.Debug("ConfigRequest from %s", event.Publisher.Id)

	poolMap := map[string]models2.DhcpPoolData{}
	leasesMap := map[string]models2.LeaseData{}
	bindingsMap := map[string]models2.DhcpBindingData{}

	if pools, err := ths.poolsCache.Get(); err == nil {
		for _, p := range pools {
			poolMap[p.Id] = p
		}
	} else {
		ths.log.Error("Can't get pools: %v", err)
	}

	if leases, err := ths.leasesCache.Get(); err == nil {
		for _, l := range leases {
			leasesMap[l.Ip] = l
		}
	} else {
		ths.log.Error("Can't get leases: %v", err)
	}

	if bindings, err := ths.bindingCache.Get(); err == nil {
		for _, b := range bindings {
			bindingsMap[b.Ip] = b
		}
	} else {
		ths.log.Error("Can't get bindings: %v", err)
	}

	return map[string]interface{}{
		"pools":    poolMap,
		"leases":   leasesMap,
		"bindings": bindingsMap,
	}
}

func (ths *EventHandler) Start() {
	if ths.startup {
		ths.startup = false
		ths.Wg.Done()
	}

	select {}
}

func NewDhcpEventHandler(
	cfg *config.Config,
	log logger.Logger,
	nasesRepo repository2.Nases,
	leasesRepo repository2.DhcpLeases,
	poolsRepo repository2.DhcpPools,
	bindingRepo repository2.DhcpBindings,
) (*EventHandler, error) {
	dhcpEvents, err := events.NewEvents(cfg, log, uuid.NewString(), Queue)
	if err != nil {
		return nil, fmt.Errorf("can't create event system: %v", err)
	}
	dhcpEvents.Run()

	broadcastEvents, err := events.NewEvents(cfg, log, uuid.NewString(), events.BroadcastEvents)
	if err != nil {
		return nil, fmt.Errorf("can't init event system: %w", err)
	}

	r := &EventHandler{
		cfg: cfg,
		log: log,
	}

	r.leasesCache, err = memory.NewDhcpLeases(log, leasesRepo, broadcastEvents)
	if err != nil {
		return nil, fmt.Errorf("can't init lease cache: %w", err)
	}

	r.nasCache, err = memory.NewNases(log, nasesRepo, broadcastEvents)
	if err != nil {
		return nil, fmt.Errorf("can't init NAS cache: %w", err)
	}

	r.poolsCache, err = memory.NewPools(log, poolsRepo, broadcastEvents)
	if err != nil {
		return nil, fmt.Errorf("can't init pools cache: %w", err)
	}

	r.bindingCache, err = memory.NewDhcpBindings(log, bindingRepo, broadcastEvents)
	if err != nil {
		return nil, fmt.Errorf("can't init bindings cache: %w", err)
	}

	err = dhcpEvents.SubscribeOnRequest(ConfigRequest, ConfigResponse, r.configRequestHandler)
	if err != nil {
		return nil, fmt.Errorf("can't subscribe on ConfigRequest: %v", err)
	}

	r.Wg.Add(1)
	r.startup = true

	return r, nil
}
