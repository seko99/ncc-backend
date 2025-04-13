package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type MapNodes struct {
	storage *psqlstorage.Storage
}

func (s *MapNodes) Create(data models.MapNodeData) error {
	r := s.storage.GetDB().Model(&models.MapNodeData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *MapNodes) Update(data models.MapNodeData) error {
	r := s.storage.GetDB().Model(&models.MapNodeData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *MapNodes) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.MapNodeData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *MapNodes) Get() ([]models.MapNodeData, error) {
	var data []models.MapNodeData

	r := s.storage.GetDB().Model(&models.MapNodeData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewMapNodes(storage *psqlstorage.Storage) *MapNodes {
	return &MapNodes{
		storage: storage,
	}
}
