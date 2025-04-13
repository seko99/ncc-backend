package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type IssueTypes struct {
	storage *psqlstorage.Storage
}

func (s *IssueTypes) Create(data models.IssueTypeData) error {
	r := s.storage.GetDB().Model(&models.IssueTypeData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *IssueTypes) Upsert(data models.IssueTypeData) error {
	r := s.storage.GetDB().Model(&models.IssueTypeData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *IssueTypes) Update(data models.IssueTypeData) error {
	r := s.storage.GetDB().Model(&models.IssueTypeData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *IssueTypes) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.IssueTypeData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *IssueTypes) Get() ([]models.IssueTypeData, error) {
	var data []models.IssueTypeData

	r := s.storage.GetDB().Model(&models.IssueTypeData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("name").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewIssueTypes(storage *psqlstorage.Storage) *IssueTypes {
	return &IssueTypes{
		storage: storage,
	}
}
