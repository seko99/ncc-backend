package memory

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNotFound = errors.New("not found")
)

type Sessions struct {
	sync.Mutex
	log      logger.Logger
	sessions []models.SessionData
	events   *events.Events
}

func (ths *Sessions) Get() ([]models.SessionData, error) {
	ths.Lock()
	defer ths.Unlock()

	result := make([]models.SessionData, len(ths.sessions))

	copy(result, ths.sessions)

	return result, nil
}

func (ths *Sessions) GetById(id string) (models.SessionData, error) {
	ths.Lock()
	defer ths.Unlock()

	for _, s := range ths.sessions {
		if s.Id == id {
			return s, nil
		}
	}

	return models.SessionData{}, ErrNotFound
}

func (ths *Sessions) GetBySessionId(sessionId string) (models.SessionData, error) {
	ths.Lock()
	defer ths.Unlock()

	for _, s := range ths.sessions {
		if s.AcctSessionId == sessionId {
			return s, nil
		}
	}

	return models.SessionData{}, ErrNotFound
}

func (ths *Sessions) GetByIP(ip string) (models.SessionData, error) {
	ths.Lock()
	defer ths.Unlock()

	for _, s := range ths.sessions {
		if s.Ip == ip {
			return s, nil
		}
	}

	return models.SessionData{}, ErrNotFound
}

func (ths *Sessions) GetByLogin(login string) ([]models.SessionData, error) {
	ths.Lock()
	defer ths.Unlock()

	var result []models.SessionData

	for _, s := range ths.sessions {
		if s.Login == login {
			result = append(result, s)
		}
	}

	if len(result) == 0 {
		return []models.SessionData{}, ErrNotFound
	}

	return result, nil
}

func (ths *Sessions) Create(data []models.SessionData) error {
	ths.Lock()
	defer ths.Unlock()

	for _, s := range data {
		ths.sessions = append(ths.sessions, s)
	}

	return nil
}

func (ths *Sessions) Update(data models.SessionData) error {
	ths.Lock()
	defer ths.Unlock()

	for idx, s := range ths.sessions {
		if s.Id == data.Id {
			ths.sessions[idx] = s
			return nil
		}
	}

	return nil
}

func (ths *Sessions) UpdateBySessionId(data models.SessionData) error {
	ths.Lock()
	defer ths.Unlock()

	for idx, s := range ths.sessions {
		if s.AcctSessionId == data.AcctSessionId {
			ths.sessions[idx] = data
			return nil
		}
	}

	return nil
}

func (ths *Sessions) GetByCustomer(id string, period repository.TimePeriod, limit ...int) ([]models.SessionData, error) {
	//TODO implement me
	panic("implement me")
}

func (ths *Sessions) Delete(id string) error {
	ths.Lock()
	defer ths.Unlock()

	var result []models.SessionData

	for _, s := range ths.sessions {
		if s.Id != id {
			result = append(result, s)
		}
	}

	ths.sessions = make([]models.SessionData, len(result))
	copy(ths.sessions, result)

	return nil
}

func (ths *Sessions) DeleteAll() error {
	ths.Lock()
	defer ths.Unlock()

	ths.sessions = []models.SessionData{}

	return nil
}

func (ths *Sessions) DeleteBySessionId(data models.SessionData) error {
	ths.Lock()
	defer ths.Unlock()

	var result []models.SessionData

	for _, s := range ths.sessions {
		if s.AcctSessionId != data.AcctSessionId {
			result = append(result, s)
		}
	}

	ths.sessions = make([]models.SessionData, len(result))
	copy(ths.sessions, result)

	return nil
}

func (ths *Sessions) onCreated(event events.Event) {
	var session models.SessionData

	b, err := json.Marshal(event.Payload)
	if err != nil {
		ths.log.Error("Can't marshal: %v", err)
		return
	}
	err = json.Unmarshal(b, &session)
	if err != nil {
		ths.log.Error("Can't unmarshal: %v", err)
		return
	}

	err = ths.Create([]models.SessionData{session})
	if err != nil {
		ths.log.Error("Can't create session: %v", err)
		return
	}

	ths.log.Debug("onCreated: %s", session.Login)
}

func (ths *Sessions) onUpdated(event events.Event) {
	var session models.SessionData

	b, err := json.Marshal(event.Payload)
	if err != nil {
		ths.log.Error("Can't marshal: %v", err)
		return
	}
	err = json.Unmarshal(b, &session)
	if err != nil {
		ths.log.Error("Can't unmarshal: %v", err)
		return
	}

	if len(session.AcctSessionId) > 0 {
		err = ths.UpdateBySessionId(session)
		if err != nil {
			ths.log.Error("Can't update session: %v", err)
			return
		}
	} else {
		err = ths.Update(session)
		if err != nil {
			ths.log.Error("Can't update session: %v", err)
			return
		}
	}

	ths.log.Debug("onUpdated: %s/%s", session.Id, session.AcctSessionId)
}

func (ths *Sessions) onAllDeleted(event events.Event) {
	err := ths.DeleteAll()
	if err != nil {
		ths.log.Error("Can't delete sessions: %v", err)
		return
	}
}

func (ths *Sessions) onDeleted(event events.Event) {
	var session models.SessionData

	b, err := json.Marshal(event.Payload)
	if err != nil {
		ths.log.Error("Can't marshal: %v", err)
		return
	}
	err = json.Unmarshal(b, &session)
	if err != nil {
		ths.log.Error("Can't unmarshal: %v", err)
		return
	}

	if len(session.AcctSessionId) > 0 {
		err = ths.DeleteBySessionId(session)
		if err != nil {
			ths.log.Error("Can't delete session: %v", err)
			return
		}
	} else {
		err = ths.Delete(session.Id)
		if err != nil {
			ths.log.Error("Can't delete session: %v", err)
			return
		}
	}

	ths.log.Debug("onDeleted: %s", session.Id)
}

func NewSessions(log logger.Logger, sessionsRepo repository.Sessions, e *events.Events) (*Sessions, error) {
	sessions := Sessions{
		log:      log,
		sessions: []models.SessionData{},
		events:   e,
	}

	if e != nil {
		err := e.SubscribeOnBroadcast(repository.SessionCreatedEvent, sessions.onCreated)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onCreated: %w", err)
		}

		err = e.SubscribeOnBroadcast(repository.SessionUpdatedEvent, sessions.onUpdated)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onUpdated: %w", err)
		}

		err = e.SubscribeOnBroadcast(repository.SessionDeletedEvent, sessions.onDeleted)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onDeleted: %w", err)
		}

		err = e.SubscribeOnBroadcast(repository.SessionAllDeletedEvent, sessions.onAllDeleted)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onAllDeleted: %w", err)
		}
	}

	persistentSessions, err := sessionsRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get sessions for cache: %w", err)
	}
	err = sessions.Create(persistentSessions)
	if err != nil {
		return nil, fmt.Errorf("can't init session cache: %w", err)
	}

	return &sessions, nil
}
