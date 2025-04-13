package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_users_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository Users

type Users interface {
	Create(data models.SecUserData) error
	Update(data models.SecUserData) error
	Delete(id string) error
	Get() ([]models.SecUserData, error)
	GetGroups() ([]models.SecGroupData, error)
}
