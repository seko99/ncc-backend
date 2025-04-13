package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_customer_groups_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository CustomerGroups

type CustomerGroups interface {
	Create(data models.CustomerGroupData) error
	Upsert(data models.CustomerGroupData) error
	Update(data models.CustomerGroupData) error
	Delete(id string) error
	Get() ([]models.CustomerGroupData, error)
}
