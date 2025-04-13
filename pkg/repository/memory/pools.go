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

type Pools struct {
	sync.Mutex
	log    logger.Logger
	pools  []models.DhcpPoolData
	events *events.Events
}

func (ths *Pools) Get() ([]models.DhcpPoolData, error) {
	ths.Lock()
	defer ths.Unlock()

	result := make([]models.DhcpPoolData, len(ths.pools))

	copy(result, ths.pools)

	return result, nil
}

func (ths *Pools) Create(data []models.DhcpPoolData) error {
	ths.Lock()
	defer ths.Unlock()

	for _, n := range data {
		ths.pools = append(ths.pools, n)
	}

	return nil
}

func (ths *Pools) Upsert(data models.DhcpPoolData) error {
	//TODO implement me
	panic("implement me")
}

func (ths *Pools) Update(data models.DhcpPoolData) error {
	ths.Lock()
	defer ths.Unlock()

	for idx, s := range ths.pools {
		if s.Id == data.Id {
			ths.pools[idx] = s
			return nil
		}
	}

	return nil
}

func (ths *Pools) Delete(id string) error {
	ths.Lock()
	defer ths.Unlock()

	var result []models.DhcpPoolData

	for _, n := range ths.pools {
		if n.Id != id {
			result = append(result, n)
		}
	}

	ths.pools = make([]models.DhcpPoolData, len(result))
	copy(ths.pools, result)

	return nil
}

func (ths *Pools) DeleteAll() error {
	ths.Lock()
	defer ths.Unlock()

	ths.pools = []models.DhcpPoolData{}

	return nil
}

func (ths *Pools) onCreated(event events.Event) {
	var pool models.DhcpPoolData

	b, err := json.Marshal(event.Payload)
	if err != nil {
		ths.log.Error("Can't marshal: %v", err)
		return
	}
	err = json.Unmarshal(b, &pool)
	if err != nil {
		ths.log.Error("Can't unmarshal: %v", err)
		return
	}

	err = ths.Create([]models.DhcpPoolData{pool})
	if err != nil {
		ths.log.Error("Can't create pool: %v", err)
		return
	}

	ths.log.Debug("onCreated: %s", pool.Id)
}

func (ths *Pools) onUpdated(event events.Event) {
	var pool models.DhcpPoolData

	b, err := json.Marshal(event.Payload)
	if err != nil {
		ths.log.Error("Can't marshal: %v", err)
		return
	}
	err = json.Unmarshal(b, &pool)
	if err != nil {
		ths.log.Error("Can't unmarshal: %v", err)
		return
	}

	err = ths.Update(pool)
	if err != nil {
		ths.log.Error("Can't update pool: %v", err)
		return
	}

	ths.log.Debug("onUpdated: %s", pool.Id)
}

func (ths *Pools) onAllDeleted(event events.Event) {
	err := ths.DeleteAll()
	if err != nil {
		ths.log.Error("Can't delete pools: %v", err)
		return
	}
}

func (ths *Pools) onDeleted(event events.Event) {
	var pool models.DhcpPoolData

	b, err := json.Marshal(event.Payload)
	if err != nil {
		ths.log.Error("Can't marshal: %v", err)
		return
	}
	err = json.Unmarshal(b, &pool)
	if err != nil {
		ths.log.Error("Can't unmarshal: %v", err)
		return
	}

	err = ths.Delete(pool.Id)
	if err != nil {
		ths.log.Error("Can't delete pool: %v", err)
		return
	}

	ths.log.Debug("onDeleted: %s", pool.Id)
}

func NewPools(log logger.Logger, poolsRepo repository.DhcpPools, e *events.Events) (*Pools, error) {
	pools := &Pools{
		log:   log,
		pools: []models.DhcpPoolData{},
	}

	if e != nil {
		err := e.SubscribeOnBroadcast(repository.NASCreatedEvent, pools.onCreated)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onCreated: %w", err)
		}

		err = e.SubscribeOnBroadcast(repository.NASUpdatedEvent, pools.onUpdated)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onUpdated: %w", err)
		}

		err = e.SubscribeOnBroadcast(repository.NASDeletedEvent, pools.onDeleted)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onDeleted: %w", err)
		}

		err = e.SubscribeOnBroadcast(repository.NASAllDeletedEvent, pools.onAllDeleted)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onAllDeleted: %w", err)
		}
	}

	persistentPools, err := poolsRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get pools for cache: %w", err)
	}
	err = pools.Create(persistentPools)
	if err != nil {
		return nil, fmt.Errorf("can't init pool cache: %w", err)
	}

	return pools, nil
}
