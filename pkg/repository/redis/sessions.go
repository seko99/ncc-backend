package redis

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	redisstorage "code.evixo.ru/ncc/ncc-backend/pkg/storage/redis"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Sessions struct {
	ctx     context.Context
	storage *redisstorage.Storage
	db      *redis.Client
}

func NewSessions(storage *redisstorage.Storage) (Sessions, error) {
	db, err := storage.GetDB(redisstorage.RedisSessionsDb)
	if err != nil {
		return Sessions{}, fmt.Errorf("can't get DB: %w", err)
	}

	return Sessions{
		ctx:     context.Background(),
		storage: storage,
		db:      db,
	}, nil
}

func (ths Sessions) Get() ([]models.SessionData, error) {
	var sessions []models.SessionData

	iter := ths.db.Scan(ths.ctx, 0, "*", 0).Iterator()
	if iter.Err() != nil {
		return nil, fmt.Errorf("can't scan: %w", iter.Err())
	}
	for iter.Next(ths.ctx) {
		var session models.SessionData
		key := iter.Val()
		cmd := ths.db.Get(ths.ctx, key)
		val, err := cmd.Result()
		if err != nil {
			return nil, fmt.Errorf("can't get value: %w", err)
		}

		err = json.Unmarshal([]byte(val), &session)
		if err != nil {
			return nil, fmt.Errorf("can't get value: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (ths Sessions) GetByIP(ip string) (models.SessionData, error) {
	iter := ths.db.Scan(ths.ctx, 0, "*", 0).Iterator()
	if iter.Err() != nil {
		return models.SessionData{}, fmt.Errorf("can't scan: %w", iter.Err())
	}
	for iter.Next(ths.ctx) {
		var s models.SessionData
		key := iter.Val()
		cmd := ths.db.Get(ths.ctx, key)
		val, err := cmd.Result()
		if err != nil {
			return models.SessionData{}, fmt.Errorf("can't get value: %w", err)
		}

		err = json.Unmarshal([]byte(val), &s)
		if err != nil {
			return models.SessionData{}, fmt.Errorf("can't get value: %w", err)
		}

		if s.Ip == ip {
			return s, nil
		}
	}

	return models.SessionData{}, fmt.Errorf("not found")
}

func (ths Sessions) GetById(id string) (models.SessionData, error) {
	//TODO implement me
	panic("implement me")
}

func (ths Sessions) GetByLogin(login string) ([]models.SessionData, error) {
	//TODO implement me
	panic("implement me")
}

func (ths Sessions) Create(data []models.SessionData) error {
	for _, d := range data {
		if len(d.CommonData.Id) == 0 {
			d.CommonData.Id = uuid.NewString()
		}

		err := ths.setSession(d.CommonData.Id, d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ths Sessions) Update(data models.SessionData) error {
	err := ths.setSession(data.CommonData.Id, data)
	if err != nil {
		return err
	}

	return nil
}

func (ths Sessions) UpdateBySessionId(data models.SessionData) error {
	//TODO implement me
	panic("implement me")
}

func (ths Sessions) GetByCustomer(id string, period repository.TimePeriod, limit ...int) ([]models.SessionData, error) {
	//TODO implement me
	panic("implement me")
}

func (ths Sessions) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

func (ths Sessions) DeleteAll() error {
	status := ths.db.FlushDB(ths.ctx)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

func (ths Sessions) DeleteBySessionId(data models.SessionData) error {
	//TODO implement me
	panic("implement me")
}

func (ths Sessions) setSession(key string, val interface{}) error {
	val, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("can't marshal value: %w", err)
	}

	status := ths.db.Set(ths.ctx, key, val, 0)
	_, err = status.Result()
	if err != nil {
		return fmt.Errorf("can't set value: %w", err)
	}

	return nil
}
