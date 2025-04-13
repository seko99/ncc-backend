package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"gorm.io/gorm/clause"
)

type NasTypes struct {
	storage *psqlstorage.Storage
}

func (s *NasTypes) Create(data models.NasTypeData) error {
	r := s.storage.GetDB().Model(&models.NasTypeData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *NasTypes) Upsert(data models.NasTypeData) error {
	r := s.storage.GetDB().Model(&models.NasTypeData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *NasTypes) Get() ([]models.NasTypeData, error) {
	var nases []models.NasTypeData

	r := s.storage.GetDB().Model(&models.NasTypeData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&nases)

	if r.Error != nil {
		return nil, r.Error
	}

	return nases, nil
}

func NewNasTypes(storage *psqlstorage.Storage) *NasTypes {
	return &NasTypes{
		storage: storage,
	}
}
