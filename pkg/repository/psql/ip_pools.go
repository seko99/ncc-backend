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

type IpPools struct {
	storage *psqlstorage.Storage
	events  *events.Events
}

func (ths *IpPools) Create(data []models.IpPoolData) error {
	r := ths.storage.GetDB().Model(&models.IpPoolData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.IpPoolCreatedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.IpPoolCreatedEvent, err)
		}
	}

	return nil
}

func (ths *IpPools) Upsert(data models.IpPoolData) error {
	r := ths.storage.GetDB().Model(&models.IpPoolData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.IpPoolCreatedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.IpPoolCreatedEvent, err)
		}
	}

	return nil
}

func (ths *IpPools) Update(data models.IpPoolData) error {
	r := ths.storage.GetDB().Model(&models.IpPoolData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.IpPoolUpdatedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.IpPoolUpdatedEvent, err)
		}
	}

	return nil
}

func (ths *IpPools) Delete(id string) error {
	r := ths.storage.GetDB().Model(&models.IpPoolData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.IpPoolDeletedEvent, &models.NasData{
			CommonData: models.CommonData{
				Id: id,
			},
		}))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.IpPoolDeletedEvent, err)
		}
	}

	return nil
}

func (ths *IpPools) DeleteAll() error {
	r := ths.storage.GetDB().Exec("DELETE FROM ?", clause.Table{Name: models.IpPoolData{}.TableName()})
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.IpPoolAllDeletedEvent, nil))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.IpPoolAllDeletedEvent, err)
		}
	}

	return nil
}

func (ths *IpPools) Get() ([]models.IpPoolData, error) {
	var data []models.IpPoolData

	r := ths.storage.GetDB().Model(&models.IpPoolData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewIpPools(storage *psqlstorage.Storage, e *events.Events) *IpPools {
	return &IpPools{
		storage: storage,
		events:  e,
	}
}
