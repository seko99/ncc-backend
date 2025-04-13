package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_issues_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository Issues

type Issues interface {
	Create(data models.IssueData) error
	Upsert(data models.IssueData) error
	Update(data models.IssueData) error
	Delete(id string) error
	DeleteAll() error
	Get() ([]models.IssueData, error)
}
