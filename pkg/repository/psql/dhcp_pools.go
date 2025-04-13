package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"fmt"
	"gorm.io/gorm/clause"
	"time"
)

type DhcpPools struct {
	storage *psqlstorage.Storage
	events  *events.Events
}

func (ths *DhcpPools) Create(data []models.DhcpPoolData) error {
	r := ths.storage.GetDB().Model(&models.DhcpPoolData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.PoolCreatedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.PoolCreatedEvent, err)
		}
	}

	return nil
}

func (ths *DhcpPools) Upsert(data models.DhcpPoolData) error {
	r := ths.storage.GetDB().Model(&models.DhcpPoolData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.PoolCreatedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.PoolCreatedEvent, err)
		}
	}

	return nil
}

func (ths *DhcpPools) Update(data models.DhcpPoolData) error {
	r := ths.storage.GetDB().Model(&models.DhcpPoolData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.PoolUpdatedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.PoolUpdatedEvent, err)
		}
	}

	return nil
}

func (ths *DhcpPools) Delete(id string) error {
	r := ths.storage.GetDB().Model(&models.DhcpPoolData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.PoolDeletedEvent, &models.NasData{
			CommonData: models.CommonData{
				Id: id,
			},
		}))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.PoolDeletedEvent, err)
		}
	}

	return nil
}

func (ths *DhcpPools) DeleteAll() error {
	r := ths.storage.GetDB().Exec("DELETE FROM ?", clause.Table{Name: models.DhcpPoolData{}.TableName()})
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.PoolAllDeletedEvent, nil))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.PoolAllDeletedEvent, err)
		}
	}

	return nil
}

func (ths *DhcpPools) Get() ([]models.DhcpPoolData, error) {
	var data []models.DhcpPoolData

	r := ths.storage.GetDB().Model(&models.DhcpPoolData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewDhcpPools(storage *psqlstorage.Storage, e *events.Events) *DhcpPools {
	return &DhcpPools{
		storage: storage,
		events:  e,
	}
}
