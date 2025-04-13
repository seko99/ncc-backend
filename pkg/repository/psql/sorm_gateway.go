package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type SormGateway struct {
	storage *psqlstorage.Storage
}

func (s *SormGateway) Create(data models.SormGatewayData) error {
	r := s.storage.GetDB().Model(&models.SormGatewayData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormGateway) Upsert(data models.SormGatewayData) error {
	r := s.storage.GetDB().Model(&models.SormGatewayData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormGateway) Update(data models.SormGatewayData) error {
	r := s.storage.GetDB().Model(&models.SormGatewayData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *SormGateway) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.SormGatewayData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormGateway) Get() ([]models.SormGatewayData, error) {
	var data []models.SormGatewayData

	r := s.storage.GetDB().Model(&models.SormGatewayData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewSormGateway(storage *psqlstorage.Storage) *SormGateway {
	return &SormGateway{
		storage: storage,
	}
}
