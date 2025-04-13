package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type Issues struct {
	storage *psqlstorage.Storage
}

func (s *Issues) Create(data models.IssueData) error {
	r := s.storage.GetDB().Model(&models.IssueData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Issues) Upsert(data models.IssueData) error {
	r := s.storage.GetDB().Model(&models.IssueData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Issues) Update(data models.IssueData) error {
	r := s.storage.GetDB().Model(&models.IssueData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *Issues) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.IssueData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Issues) DeleteAll() error {
	r := s.storage.GetDB().Exec("DELETE FROM ?", clause.Table{Name: models.IssueData{}.TableName()})
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *Issues) Get() ([]models.IssueData, error) {
	var data []models.IssueData

	r := s.storage.GetDB().Model(&models.IssueData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("create_ts").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewIssues(storage *psqlstorage.Storage) *Issues {
	return &Issues{
		storage: storage,
	}
}
