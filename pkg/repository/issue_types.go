package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_issue_types_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository IssueTypes

type IssueTypes interface {
	Create(data models.IssueTypeData) error
	Upsert(data models.IssueTypeData) error
	Update(data models.IssueTypeData) error
	Delete(id string) error
	Get() ([]models.IssueTypeData, error)
}
