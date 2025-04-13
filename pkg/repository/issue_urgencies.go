package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_issue_urgencies_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository IssueUrgencies

type IssueUrgencies interface {
	Create(data models.IssueUrgencyData) error
	Upsert(data models.IssueUrgencyData) error
	Update(data models.IssueUrgencyData) error
	Delete(id string) error
	Get() ([]models.IssueUrgencyData, error)
}
