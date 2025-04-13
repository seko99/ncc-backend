package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type IssueActions struct {
	storage *psqlstorage.Storage
}

func (s *IssueActions) Create(data models.IssueActionData) error {
	r := s.storage.GetDB().Model(&models.IssueActionData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *IssueActions) Upsert(data models.IssueActionData) error {
	r := s.storage.GetDB().Model(&models.IssueActionData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *IssueActions) Update(data models.IssueActionData) error {
	r := s.storage.GetDB().Model(&models.IssueActionData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *IssueActions) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.IssueActionData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *IssueActions) Get() ([]models.IssueActionData, error) {
	var data []models.IssueActionData

	r := s.storage.GetDB().Model(&models.IssueActionData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewIssueActions(storage *psqlstorage.Storage) *IssueActions {
	return &IssueActions{
		storage: storage,
	}
}
