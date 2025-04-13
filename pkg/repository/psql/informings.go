package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"gorm.io/gorm/clause"
	"time"
)

type Informings struct {
	storage *psqlstorage.Storage
}

func (s *Informings) Get() ([]models.InformingData, error) {
	var informings []models.InformingData
	r := s.storage.GetDB().Model(models.InformingData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Find(&informings)
	if r.Error != nil {
		return nil, r.Error
	}
	return informings, nil
}

func (s *Informings) GetEnabled() ([]models.InformingData, error) {
	var informings []models.InformingData
	r := s.storage.GetDB().Model(models.InformingData{}).
		Preload(clause.Associations).
		Preload("Conditions", "delete_ts is null").
		Where("delete_ts is null").
		Where("state = ?", models.InformingStateEnabled).
		Find(&informings)
	if r.Error != nil {
		return nil, r.Error
	}
	return informings, nil
}

func (s *Informings) Create(data models.InformingData) error {
	r := s.storage.GetDB().Model(&models.InformingData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Informings) SetState(data models.InformingData, state int) error {
	r := s.storage.GetDB().Model(&models.InformingData{}).
		Where("id = ?", data.Id).
		Update("state", state)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Informings) SetStart(data models.InformingData, start time.Time) error {
	r := s.storage.GetDB().Model(&models.InformingData{}).
		Where("id = ?", data.Id).
		Update("start", models.NewNullTime(start))
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func NewInformings(storage *psqlstorage.Storage) *Informings {
	return &Informings{
		storage: storage,
	}
}
