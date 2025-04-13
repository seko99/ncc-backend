package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"time"
)

//go:generate mockgen -destination=mocks/mock_customers_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository Customers

const (
	CustomerCreatedEvent    = "customerCreated"
	CustomerUpdatedEvent    = "customerUpdated"
	CustomerDeletedEvent    = "customerDeleted"
	CustomerAllDeletedEvent = "customerAllDeleted"
)

type Customers interface {
	Create(data models.CustomerData) error
	Update(data models.CustomerData) error
	DeleteAll() error

	Get(limit ...int) ([]models.CustomerData, error)
	GetById(id string) (*models.CustomerData, error)
	GetByUid(uid string) (*models.CustomerData, error)
	GetByLogin(login string) (*models.CustomerData, error)
	GetByState(state int) ([]models.CustomerData, error)
	GetByGroup(id string) ([]models.CustomerData, error)

	GetByFeeAmountAndSessions(amount float64, sdate string) ([]models.CustomerData, error)

	GetGroupByName(name string) (*models.CustomerGroupData, error)
	IsBlocked(state int) bool
	SetDeposit(id string, deposit float64) error
	SetCredit(id string, credit float64) error
	SetCreditExpire(id string, creditExpire time.Time) error
	SetCreditDaysLeft(id string, creditDaysLeft int) error
	SetState(id string, state int) error
	SetServiceInternetState(id string, state int) error
	SetFlag(customer models.CustomerData, flag models.CustomerFlagData) error
}
