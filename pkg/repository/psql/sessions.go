package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"fmt"
	"gorm.io/gorm/clause"
	"time"
)

type Sessions struct {
	storage *psqlstorage.Storage
	events  *events.Events
}

func (ths *Sessions) Create(sessions []models2.SessionData) error {
	r := ths.storage.GetDB().Model(&models2.SessionData{}).
		Create(sessions)

	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		for _, s := range sessions {
			err := ths.events.PublishEvent(events.NewEvent(repository.SessionCreatedEvent, s))
			if err != nil {
				return fmt.Errorf("can't publish event %s: %w", repository.SessionCreatedEvent, err)
			}
		}
	}

	return nil
}

func (ths *Sessions) Get() ([]models2.SessionData, error) {
	var sessions []models2.SessionData

	r := ths.storage.GetDB().Model(&models2.SessionData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("login").
		Find(&sessions)

	if r.Error != nil {
		return nil, r.Error
	}

	return sessions, nil
}

func (ths *Sessions) GetById(id string) (models2.SessionData, error) {
	var session models2.SessionData

	r := ths.storage.GetDB().Model(&models2.SessionData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Order("login").
		First(&session)

	if r.Error != nil {
		return models2.SessionData{}, r.Error
	}

	return session, nil
}

func (ths *Sessions) GetBySessionId(sessionId string) (models2.SessionData, error) {
	var session models2.SessionData

	r := ths.storage.GetDB().Model(&models2.SessionData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("acct_session_id = @session_id", sql.Named("session_id", sessionId)).
		Order("login").
		First(&session)

	if r.Error != nil {
		return models2.SessionData{}, r.Error
	}

	return session, nil
}

func (ths *Sessions) GetByIP(ip string) (models2.SessionData, error) {
	var session models2.SessionData

	r := ths.storage.GetDB().Model(&models2.SessionData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("ip = @ip", sql.Named("ip", ip)).
		Order("login").
		First(&session)

	if r.Error != nil {
		return models2.SessionData{}, r.Error
	}

	return session, nil
}

func (ths *Sessions) GetByLogin(login string) ([]models2.SessionData, error) {
	var sessions []models2.SessionData

	r := ths.storage.GetDB().Model(&models2.SessionData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("login = @login", sql.Named("login", login)).
		Order("login").
		First(&sessions)

	if r.Error != nil {
		return nil, r.Error
	}

	return sessions, nil
}

func (ths *Sessions) Update(data models2.SessionData) error {
	r := ths.storage.GetDB().Model(&models2.SessionData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.SessionUpdatedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.SessionUpdatedEvent, err)
		}
	}

	return nil
}

func (ths *Sessions) UpdateBySessionId(data models2.SessionData) error {
	r := ths.storage.GetDB().Model(&models2.SessionData{}).
		Where("acct_session_id = @id", sql.Named("id", data.AcctSessionId)).
		Updates(map[string]interface{}{
			"update_ts":  time.Now(),
			"duration":   data.Duration,
			"last_alive": data.LastAlive,
			"octets_in":  data.OctetsIn,
			"octets_out": data.OctetsOut,
		})
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.SessionUpdatedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.SessionUpdatedEvent, err)
		}
	}

	return nil
}

func (ths *Sessions) Delete(id string) error {
	r := ths.storage.GetDB().Model(&models2.SessionData{}).
		Where("id = @id", sql.Named("id", id)).
		Delete(models2.SessionData{
			CommonData: models2.CommonData{
				Id: id,
			},
		})
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.SessionDeletedEvent, &models2.SessionData{
			CommonData: models2.CommonData{
				Id: id,
			},
		}))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.SessionDeletedEvent, err)
		}
	}

	return nil
}

func (ths *Sessions) DeleteBySessionId(data models2.SessionData) error {
	r := ths.storage.GetDB().
		Where("acct_session_id = @id", sql.Named("id", data.AcctSessionId)).
		Delete(&models2.SessionData{})
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.SessionDeletedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.SessionDeletedEvent, err)
		}
	}

	return nil
}

func (ths *Sessions) DeleteAll() error {
	r := ths.storage.GetDB().Exec("DELETE FROM ?", clause.Table{Name: models2.SessionData{}.TableName()})
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.SessionAllDeletedEvent, nil))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.SessionAllDeletedEvent, err)
		}
	}

	return nil
}

func (ths *Sessions) GetByCustomer(id string, period repository.TimePeriod, limit ...int) ([]models2.SessionData, error) {
	var sessions []models2.SessionData

	periodClause := repository.PeriodClause("start_time", period)
	r := ths.storage.GetDB().Model(&models2.SessionData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("customer_id = @id", sql.Named("id", id)).
		Where(periodClause).
		Order("start_time")

	if len(limit) > 0 {
		r = r.Limit(limit[0])
	}

	r = r.Find(&sessions)

	if r.Error != nil {
		return nil, r.Error
	}

	return sessions, nil
}

func NewSessions(storage *psqlstorage.Storage, e *events.Events) *Sessions {
	sessions := &Sessions{
		storage: storage,
		events:  e,
	}
	return sessions
}
