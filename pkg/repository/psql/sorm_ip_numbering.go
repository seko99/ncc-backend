package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type SormIpNumbering struct {
	storage *psqlstorage.Storage
}

func (s *SormIpNumbering) Create(data models.SormIpNumberingData) error {
	r := s.storage.GetDB().Model(&models.SormIpNumberingData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormIpNumbering) Upsert(data models.SormIpNumberingData) error {
	r := s.storage.GetDB().Model(&models.SormIpNumberingData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormIpNumbering) Update(data models.SormIpNumberingData) error {
	r := s.storage.GetDB().Model(&models.SormIpNumberingData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *SormIpNumbering) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.SormIpNumberingData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SormIpNumbering) Get() ([]models.SormIpNumberingData, error) {
	var data []models.SormIpNumberingData

	r := s.storage.GetDB().Model(&models.SormIpNumberingData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewSormIpNumbering(storage *psqlstorage.Storage) *SormIpNumbering {
	return &SormIpNumbering{
		storage: storage,
	}
}
