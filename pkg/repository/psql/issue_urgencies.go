package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type IssueUrgencies struct {
	storage *psqlstorage.Storage
}

func (s *IssueUrgencies) Create(data models.IssueUrgencyData) error {
	r := s.storage.GetDB().Model(&models.IssueUrgencyData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *IssueUrgencies) Upsert(data models.IssueUrgencyData) error {
	r := s.storage.GetDB().Model(&models.IssueUrgencyData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *IssueUrgencies) Update(data models.IssueUrgencyData) error {
	r := s.storage.GetDB().Model(&models.IssueUrgencyData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *IssueUrgencies) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.IssueUrgencyData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *IssueUrgencies) Get() ([]models.IssueUrgencyData, error) {
	var data []models.IssueUrgencyData

	r := s.storage.GetDB().Model(&models.IssueUrgencyData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewIssueUrgencies(storage *psqlstorage.Storage) *IssueUrgencies {
	return &IssueUrgencies{
		storage: storage,
	}
}
