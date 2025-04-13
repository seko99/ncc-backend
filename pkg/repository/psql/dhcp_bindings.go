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

type DhcpBindings struct {
	storage *psqlstorage.Storage
	events  *events.Events
}

func (ths *DhcpBindings) Create(data []models.DhcpBindingData) error {
	r := ths.storage.GetDB().Model(&models.DhcpBindingData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		for _, b := range data {
			err := ths.events.PublishEvent(events.NewEvent(repository.BindingCreatedEvent, b))
			if err != nil {
				return fmt.Errorf("can't publish event %s: %w", repository.BindingCreatedEvent, err)
			}
		}
	}

	return nil
}

func (ths *DhcpBindings) Update(data models.DhcpBindingData) error {
	r := ths.storage.GetDB().Model(&models.DhcpBindingData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.BindingUpdatedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.BindingUpdatedEvent, err)
		}
	}

	return nil
}

func (ths *DhcpBindings) Delete(id string) error {
	r := ths.storage.GetDB().Model(&models.DhcpBindingData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.BindingDeletedEvent, &models.DhcpBindingData{
			CommonData: models.CommonData{
				Id: id,
			},
		}))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.BindingDeletedEvent, err)
		}
	}

	return nil
}

func (ths *DhcpBindings) DeleteAll() error {
	r := ths.storage.GetDB().Exec("DELETE FROM ?", clause.Table{Name: models.DhcpBindingData{}.TableName()})
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.BindingAllDeletedEvent, nil))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.BindingAllDeletedEvent, err)
		}
	}

	return nil
}

func (ths *DhcpBindings) Get() ([]models.DhcpBindingData, error) {
	var data []models.DhcpBindingData

	r := ths.storage.GetDB().Model(&models.DhcpBindingData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewDhcpBindings(storage *psqlstorage.Storage, e *events.Events) *DhcpBindings {
	bindings := &DhcpBindings{
		storage: storage,
		events:  e,
	}
	return bindings
}
