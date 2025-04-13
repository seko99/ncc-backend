package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_issue_actions_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository IssueActions

type IssueActions interface {
	Create(data models.IssueActionData) error
	Upsert(data models.IssueActionData) error
	Update(data models.IssueActionData) error
	Delete(id string) error
	Get() ([]models.IssueActionData, error)
}
