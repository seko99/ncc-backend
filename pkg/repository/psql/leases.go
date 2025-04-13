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

type DhcpLeases struct {
	storage *psqlstorage.Storage
	events  *events.Events
}

func (ths *DhcpLeases) Create(data models.LeaseData) error {
	r := ths.storage.GetDB().Model(&models.LeaseData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.LeaseCreatedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.LeaseCreatedEvent, err)
		}
	}

	return nil
}

func (ths *DhcpLeases) Update(data models.LeaseData) error {
	r := ths.storage.GetDB().Model(&models.LeaseData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.LeaseUpdatedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.LeaseUpdatedEvent, err)
		}
	}

	return nil
}

func (ths *DhcpLeases) Get() ([]models.LeaseData, error) {
	var leases []models.LeaseData

	r := ths.storage.GetDB().Model(&models.LeaseData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("start").
		Find(&leases)

	if r.Error != nil {
		return nil, r.Error
	}

	return leases, nil
}

func (ths *DhcpLeases) GetByIP(ip string) (models.LeaseData, error) {
	var lease models.LeaseData

	r := ths.storage.GetDB().Model(&models.LeaseData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("ip = @ip", sql.Named("ip", ip)).
		Order("start").
		First(&lease)

	if r.Error != nil {
		return models.LeaseData{}, r.Error
	}

	return lease, nil
}

func (ths *DhcpLeases) DeleteAll() error {
	r := ths.storage.GetDB().Exec("DELETE FROM ?", clause.Table{Name: models.LeaseData{}.TableName()})
	if r.Error != nil {
		return r.Error
	}

	if ths.events != nil {
		err := ths.events.PublishEvent(events.NewEvent(repository.LeaseAllDeletedEvent, nil))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.LeaseAllDeletedEvent, err)
		}
	}

	return nil
}

func NewDhcpLeases(storage *psqlstorage.Storage, e *events.Events) *DhcpLeases {
	leases := &DhcpLeases{
		storage: storage,
		events:  e,
	}
	return leases
}
