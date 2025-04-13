package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"fmt"
	"gorm.io/gorm/clause"
)

type Nases struct {
	storage *psqlstorage.Storage
	events  *events.Events
}

func (ths *Nases) Create(data []models.NasData) error {
	r := ths.storage.GetDB().Model(&models.NasData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.NASCreatedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.NASCreatedEvent, err)
		}
	}

	return nil
}

func (ths *Nases) Upsert(data models.NasData) error {
	r := ths.storage.GetDB().Model(&models.NasData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.NASCreatedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.NASCreatedEvent, err)
		}
	}

	return nil
}

func (ths *Nases) Update(data models.NasData) error {
	r := ths.storage.GetDB().Model(&models.NasData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.NASUpdatedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.NASUpdatedEvent, err)
		}
	}

	return nil
}

func (ths *Nases) Delete(id string) error {
	r := ths.storage.GetDB().Model(&models.NasData{}).
		Where("id = @id", sql.Named("id", id)).
		Delete(models.NasData{
			CommonData: models.CommonData{
				Id: id,
			},
		})
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.NASDeletedEvent, &models.NasData{
			CommonData: models.CommonData{
				Id: id,
			},
		}))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.NASDeletedEvent, err)
		}
	}

	return nil
}

func (ths *Nases) DeleteAll() error {
	r := ths.storage.GetDB().Exec("DELETE FROM ?", clause.Table{Name: models.NasData{}.TableName()})
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.NASAllDeletedEvent, nil))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.NASAllDeletedEvent, err)
		}
	}

	return nil
}

func (ths *Nases) Get() ([]models.NasData, error) {
	var nases []models.NasData

	r := ths.storage.GetDB().Model(&models.NasData{}).
		Preload("NasType.NasAttributes.Attribute.Vendor").
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&nases)

	if r.Error != nil {
		return nil, r.Error
	}

	return nases, nil
}

func (ths *Nases) GetByIP(ip string) (models.NasData, error) {
	var nas models.NasData

	r := ths.storage.GetDB().Model(&models.NasData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("ip = @ip", sql.Named("ip", ip)).
		Find(&nas)

	if r.Error != nil {
		return models.NasData{}, r.Error
	}

	return nas, nil
}

func NewNases(storage *psqlstorage.Storage, e *events.Events) *Nases {
	return &Nases{
		storage: storage,
		events:  e,
	}
}
