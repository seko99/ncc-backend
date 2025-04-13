package memory

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"encoding/json"
	"fmt"
	"sync"
)

type DhcpLeases struct {
	sync.Mutex
	log    logger.Logger
	leases []models.LeaseData
	events *events.Events
}

func (ths *DhcpLeases) Get() ([]models.LeaseData, error) {
	ths.Lock()
	defer ths.Unlock()

	result := make([]models.LeaseData, len(ths.leases))

	copy(result, ths.leases)

	return result, nil
}

func (ths *DhcpLeases) GetByIP(ip string) (models.LeaseData, error) {
	ths.Lock()
	defer ths.Unlock()

	for _, l := range ths.leases {
		if l.Ip == ip {
			return l, nil
		}
	}

	return models.LeaseData{}, ErrNotFound
}

func (ths *DhcpLeases) Create(data models.LeaseData) error {
	ths.Lock()
	defer ths.Unlock()

	ths.leases = append(ths.leases, data)

	return nil
}

func (ths *DhcpLeases) Update(data models.LeaseData) error {
	//TODO implement me
	panic("implement me")
}

func (ths *DhcpLeases) DeleteAll() error {
	ths.Lock()
	defer ths.Unlock()

	ths.leases = []models.LeaseData{}

	return nil
}

func (ths *DhcpLeases) onCreated(event events.Event) {
	var lease models.LeaseData

	b, err := json.Marshal(event.Payload)
	if err != nil {
		ths.log.Error("Can't marshal: %v", err)
		return
	}
	err = json.Unmarshal(b, &lease)
	if err != nil {
		ths.log.Error("Can't unmarshal: %v", err)
		return
	}

	err = ths.Create(lease)
	if err != nil {
		ths.log.Error("Can't create lease: %v", err)
		return
	}

	ths.log.Debug("onCreated: %s", lease.Ip)
}

func (ths *DhcpLeases) onUpdated(event events.Event) {
	var lease models.LeaseData

	b, err := json.Marshal(event.Payload)
	if err != nil {
		ths.log.Error("Can't marshal: %v", err)
		return
	}
	err = json.Unmarshal(b, &lease)
	if err != nil {
		ths.log.Error("Can't unmarshal: %v", err)
		return
	}

	err = ths.Update(lease)
	if err != nil {
		ths.log.Error("Can't update lease: %v", err)
		return
	}

	ths.log.Debug("onUpdated: %s", lease.Id)
}

func (ths *DhcpLeases) onAllDeleted(event events.Event) {
}

func (ths *DhcpLeases) onDeleted(event events.Event) {
	var lease models.LeaseData

	b, err := json.Marshal(event.Payload)
	if err != nil {
		ths.log.Error("Can't marshal: %v", err)
		return
	}
	err = json.Unmarshal(b, &lease)
	if err != nil {
		ths.log.Error("Can't unmarshal: %v", err)
		return
	}

	ths.log.Debug("onDeleted: %s", lease.Id)
}

func NewDhcpLeases(log logger.Logger, leasesRepo repository.DhcpLeases, e *events.Events) (*DhcpLeases, error) {
	leases := DhcpLeases{
		log:    log,
		leases: []models.LeaseData{},
	}

	if e != nil {
		err := e.SubscribeOnBroadcast(repository.LeaseCreatedEvent, leases.onCreated)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onCreated: %w", err)
		}

		err = e.SubscribeOnBroadcast(repository.LeaseUpdatedEvent, leases.onUpdated)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onUpdated: %w", err)
		}

		err = e.SubscribeOnBroadcast(repository.LeaseDeletedEvent, leases.onDeleted)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onDeleted: %w", err)
		}

		err = e.SubscribeOnBroadcast(repository.LeaseAllDeletedEvent, leases.onAllDeleted)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onAllDeleted: %w", err)
		}
	}

	persistentLeases, err := leasesRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get leases for cache: %w", err)
	}
	for _, l := range persistentLeases {
		err = leases.Create(l)
		if err != nil {
			return nil, fmt.Errorf("can't init lease cache: %w", err)
		}
	}

	return &leases, nil
}
