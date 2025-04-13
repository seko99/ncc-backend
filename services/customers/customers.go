package customers

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
)

type Customers struct {
	repo repository.Customers
}

func (s *Customers) GetByLogin(login string) (*models.CustomerData, error) {
	return s.repo.GetByLogin(login)
}

func NewCustomers(repo repository.Customers) *Customers {
	return &Customers{
		repo: repo,
	}
}
