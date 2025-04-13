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

type DhcpBindings struct {
	sync.Mutex
	log      logger.Logger
	bindings []models.DhcpBindingData
	events   *events.Events
}

func (ths *DhcpBindings) Create(data []models.DhcpBindingData) error {
	ths.Lock()
	defer ths.Unlock()

	for _, b := range data {
		ths.bindings = append(ths.bindings, b)
	}

	return nil
}

func (ths *DhcpBindings) Update(data models.DhcpBindingData) error {
	ths.Lock()
	defer ths.Unlock()

	for idx, b := range ths.bindings {
		if b.Id == data.Id {
			ths.bindings[idx] = b
			return nil
		}
	}

	return nil
}

func (ths *DhcpBindings) Delete(id string) error {
	ths.Lock()
	defer ths.Unlock()

	var result []models.DhcpBindingData

	for _, b := range ths.bindings {
		if b.Id != id {
			result = append(result, b)
		}
	}

	ths.bindings = make([]models.DhcpBindingData, len(result))
	copy(ths.bindings, result)

	return nil
}

func (ths *DhcpBindings) DeleteAll() error {
	ths.Lock()
	defer ths.Unlock()

	ths.bindings = []models.DhcpBindingData{}

	return nil
}

func (ths *DhcpBindings) Get() ([]models.DhcpBindingData, error) {
	ths.Lock()
	defer ths.Unlock()

	result := make([]models.DhcpBindingData, len(ths.bindings))

	copy(result, ths.bindings)

	return result, nil
}

func (ths *DhcpBindings) onCreated(event events.Event) {
	var binding models.DhcpBindingData

	b, err := json.Marshal(event.Payload)
	if err != nil {
		ths.log.Error("Can't marshal: %v", err)
		return
	}
	err = json.Unmarshal(b, &binding)
	if err != nil {
		ths.log.Error("Can't unmarshal: %v", err)
		return
	}

	err = ths.Create([]models.DhcpBindingData{binding})
	if err != nil {
		ths.log.Error("Can't create binding: %v", err)
		return
	}

	ths.log.Debug("onCreated: %s", binding.Id)
}

func (ths *DhcpBindings) onUpdated(event events.Event) {
	var binding models.DhcpBindingData

	b, err := json.Marshal(event.Payload)
	if err != nil {
		ths.log.Error("Can't marshal: %v", err)
		return
	}
	err = json.Unmarshal(b, &binding)
	if err != nil {
		ths.log.Error("Can't unmarshal: %v", err)
		return
	}

	err = ths.Update(binding)
	if err != nil {
		ths.log.Error("Can't update binding: %v", err)
		return
	}

	ths.log.Debug("onUpdated: %s", binding.Id)
}

func (ths *DhcpBindings) onAllDeleted(event events.Event) {
	err := ths.DeleteAll()
	if err != nil {
		ths.log.Error("Can't delete bindings: %v", err)
		return
	}
}

func (ths *DhcpBindings) onDeleted(event events.Event) {
	var binding models.DhcpBindingData

	b, err := json.Marshal(event.Payload)
	if err != nil {
		ths.log.Error("Can't marshal: %v", err)
		return
	}
	err = json.Unmarshal(b, &binding)
	if err != nil {
		ths.log.Error("Can't unmarshal: %v", err)
		return
	}

	err = ths.Delete(binding.Id)
	if err != nil {
		ths.log.Error("Can't delete binding: %v", err)
		return
	}

	ths.log.Debug("onDeleted: %s", binding.Id)
}

func NewDhcpBindings(log logger.Logger, bindingsRepo repository.DhcpBindings, e *events.Events) (*DhcpBindings, error) {
	bindings := DhcpBindings{
		log:      log,
		bindings: []models.DhcpBindingData{},
		events:   e,
	}

	if e != nil {
		err := e.SubscribeOnBroadcast(repository.BindingCreatedEvent, bindings.onCreated)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onCreated: %w", err)
		}

		err = e.SubscribeOnBroadcast(repository.BindingUpdatedEvent, bindings.onUpdated)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onUpdated: %w", err)
		}

		err = e.SubscribeOnBroadcast(repository.BindingDeletedEvent, bindings.onDeleted)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onDeleted: %w", err)
		}

		err = e.SubscribeOnBroadcast(repository.BindingAllDeletedEvent, bindings.onAllDeleted)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onAllDeleted: %w", err)
		}
	}

	persistentBindings, err := bindingsRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get bindings for cache: %w", err)
	}
	err = bindings.Create(persistentBindings)
	if err != nil {
		return nil, fmt.Errorf("can't init bindings cache: %w", err)
	}

	return &bindings, nil
}
