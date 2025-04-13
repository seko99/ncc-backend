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

type Nases struct {
	sync.Mutex
	log    logger.Logger
	nases  []models.NasData
	events *events.Events
}

func (ths *Nases) Get() ([]models.NasData, error) {
	ths.Lock()
	defer ths.Unlock()

	result := make([]models.NasData, len(ths.nases))

	copy(result, ths.nases)

	return result, nil
}

func (ths *Nases) GetByIP(ip string) (models.NasData, error) {
	ths.Lock()
	defer ths.Unlock()

	for _, n := range ths.nases {
		if n.Ip == ip {
			return n, nil
		}
	}

	return models.NasData{}, ErrNotFound
}

func (ths *Nases) Create(data []models.NasData) error {
	ths.Lock()
	defer ths.Unlock()

	for _, n := range data {
		ths.nases = append(ths.nases, n)
	}

	return nil
}

func (ths *Nases) Upsert(data models.NasData) error {
	//TODO implement me
	panic("implement me")
}

func (ths *Nases) Update(data models.NasData) error {
	ths.Lock()
	defer ths.Unlock()

	for idx, s := range ths.nases {
		if s.Id == data.Id {
			ths.nases[idx] = s
			return nil
		}
	}

	return nil
}

func (ths *Nases) Delete(id string) error {
	ths.Lock()
	defer ths.Unlock()

	var result []models.NasData

	for _, n := range ths.nases {
		if n.Id != id {
			result = append(result, n)
		}
	}

	ths.nases = make([]models.NasData, len(result))
	copy(ths.nases, result)

	return nil
}

func (ths *Nases) DeleteAll() error {
	ths.Lock()
	defer ths.Unlock()

	ths.nases = []models.NasData{}

	return nil
}

func (ths *Nases) onCreated(event events.Event) {
	var nas models.NasData

	b, err := json.Marshal(event.Payload)
	if err != nil {
		ths.log.Error("Can't marshal: %v", err)
		return
	}
	err = json.Unmarshal(b, &nas)
	if err != nil {
		ths.log.Error("Can't unmarshal: %v", err)
		return
	}

	err = ths.Create([]models.NasData{nas})
	if err != nil {
		ths.log.Error("Can't create NAS: %v", err)
		return
	}

	ths.log.Debug("onCreated: %s", nas.Id)
}

func (ths *Nases) onUpdated(event events.Event) {
	var nas models.NasData

	b, err := json.Marshal(event.Payload)
	if err != nil {
		ths.log.Error("Can't marshal: %v", err)
		return
	}
	err = json.Unmarshal(b, &nas)
	if err != nil {
		ths.log.Error("Can't unmarshal: %v", err)
		return
	}

	err = ths.Update(nas)
	if err != nil {
		ths.log.Error("Can't update NAS: %v", err)
		return
	}

	ths.log.Debug("onUpdated: %s", nas.Id)
}

func (ths *Nases) onAllDeleted(event events.Event) {
	err := ths.DeleteAll()
	if err != nil {
		ths.log.Error("Can't delete NASes: %v", err)
		return
	}
}

func (ths *Nases) onDeleted(event events.Event) {
	var nas models.NasData

	b, err := json.Marshal(event.Payload)
	if err != nil {
		ths.log.Error("Can't marshal: %v", err)
		return
	}
	err = json.Unmarshal(b, &nas)
	if err != nil {
		ths.log.Error("Can't unmarshal: %v", err)
		return
	}

	err = ths.Delete(nas.Id)
	if err != nil {
		ths.log.Error("Can't delete NAS: %v", err)
		return
	}

	ths.log.Debug("onDeleted: %s", nas.Id)
}

func NewNases(log logger.Logger, nasesRepo repository.Nases, e *events.Events) (*Nases, error) {
	nases := &Nases{
		log:   log,
		nases: []models.NasData{},
	}

	if e != nil {
		err := e.SubscribeOnBroadcast(repository.NASCreatedEvent, nases.onCreated)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onCreated: %w", err)
		}

		err = e.SubscribeOnBroadcast(repository.NASUpdatedEvent, nases.onUpdated)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onUpdated: %w", err)
		}

		err = e.SubscribeOnBroadcast(repository.NASDeletedEvent, nases.onDeleted)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onDeleted: %w", err)
		}

		err = e.SubscribeOnBroadcast(repository.NASAllDeletedEvent, nases.onAllDeleted)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onAllDeleted: %w", err)
		}
	}

	persistentNases, err := nasesRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get NASes for cache: %w", err)
	}
	err = nases.Create(persistentNases)
	if err != nil {
		return nil, fmt.Errorf("can't init NAS cache: %w", err)
	}

	return nases, nil
}
