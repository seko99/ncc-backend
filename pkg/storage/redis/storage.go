package redisstorage

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

const (
	RedisCustomersDb = 0
	RedisLeasesDb    = 1
	RedisSessionsDb  = 2
)

type Storage struct {
	cfg  *config.Config
	log  logger.Logger
	rdbs []*redis.Client
}

func (ths *Storage) GetDB(id int) (*redis.Client, error) {
	if len(ths.rdbs) < id {
		return nil, fmt.Errorf("no such DB: %d", id)
	}

	return ths.rdbs[id], nil
}

func (ths *Storage) Connect() error {
	ths.rdbs = []*redis.Client{
		RedisCustomersDb: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", ths.cfg.Redis.Host, ths.cfg.Redis.Port),
			Password: ths.cfg.Redis.Password,
			DB:       RedisCustomersDb,
		}),
		RedisLeasesDb: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", ths.cfg.Redis.Host, ths.cfg.Redis.Port),
			Password: ths.cfg.Redis.Password,
			DB:       RedisLeasesDb,
		}),
		RedisSessionsDb: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", ths.cfg.Redis.Host, ths.cfg.Redis.Port),
			Password: ths.cfg.Redis.Password,
			DB:       RedisSessionsDb,
		}),
	}

	for _, r := range ths.rdbs {
		err := r.Ping(context.Background()).Err()
		if err != nil {
			return fmt.Errorf("can't connect to redis: %v", err)
		}
	}

	return nil
}

func NewRedis(cfg *config.Config, log logger.Logger) *Storage {
	return &Storage{
		cfg: cfg,
		log: log,
	}
}
