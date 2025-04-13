package memory

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type Customers struct {
	sync.Mutex
	log       logger.Logger
	customers []models.CustomerData
	repo      repository.Customers
	events    *events.Events
	cache     *cache.Cache
}

func (ths *Customers) Create(data models.CustomerData) error {
	ths.Lock()
	defer ths.Unlock()

	ths.customers = append(ths.customers, data)

	return nil
}

func (ths *Customers) Update(data models.CustomerData) error {
	//TODO implement me
	panic("implement me")
}

func (ths *Customers) Get(limit ...int) ([]models.CustomerData, error) {
	ths.Lock()
	defer ths.Unlock()

	result := make([]models.CustomerData, len(ths.customers))

	copy(result, ths.customers)

	return result, nil
}

func (ths *Customers) GetById(id string) (*models.CustomerData, error) {
	//TODO implement me
	panic("implement me")
}

func (ths *Customers) GetByUid(uid string) (*models.CustomerData, error) {
	//TODO implement me
	panic("implement me")
}

func (ths *Customers) GetByLogin(login string) (*models.CustomerData, error) {
	ths.Lock()
	defer ths.Unlock()

	for _, c := range ths.customers {
		if c.Login == login {
			return &c, nil
		}
	}

	return &models.CustomerData{}, ErrNotFound
}

func (ths *Customers) GetByState(state int) ([]models.CustomerData, error) {
	//TODO implement me
	panic("implement me")
}

func (ths *Customers) GetByGroup(id string) ([]models.CustomerData, error) {
	//TODO implement me
	panic("implement me")
}

func (ths *Customers) GetByFeeAmountAndSessions(amount float64, sdate string) ([]models.CustomerData, error) {
	//TODO implement me
	panic("implement me")
}

func (ths *Customers) GetGroupByName(name string) (*models.CustomerGroupData, error) {
	//TODO implement me
	panic("implement me")
}

func (ths *Customers) IsBlocked(state int) bool {
	//TODO implement me
	panic("implement me")
}

func (ths *Customers) SetDeposit(id string, deposit float64) error {
	//TODO implement me
	panic("implement me")
}

func (ths *Customers) SetCredit(id string, credit float64) error {
	//TODO implement me
	panic("implement me")
}

func (ths *Customers) SetCreditExpire(id string, creditExpire time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (ths *Customers) SetCreditDaysLeft(id string, creditDaysLeft int) error {
	//TODO implement me
	panic("implement me")
}

func (ths *Customers) SetState(id string, state int) error {
	//TODO implement me
	panic("implement me")
}

func (ths *Customers) SetServiceInternetState(id string, state int) error {
	//TODO implement me
	panic("implement me")
}

func (ths *Customers) SetFlag(customer models.CustomerData, flag models.CustomerFlagData) error {
	//TODO implement me
	panic("implement me")
}

func (ths *Customers) DeleteAll() error {
	ths.Lock()
	defer ths.Unlock()

	ths.customers = []models.CustomerData{}

	return nil
}

func (ths *Customers) onCreated(event events.Event) {
	var customer models.CustomerData

	b, err := json.Marshal(event.Payload)
	if err != nil {
		ths.log.Error("Can't marshal: %v", err)
		return
	}
	err = json.Unmarshal(b, &customer)
	if err != nil {
		ths.log.Error("Can't unmarshal: %v", err)
		return
	}

	err = ths.Create(customer)
	if err != nil {
		ths.log.Error("Can't create customers: %v", err)
		return
	}
}

func (ths *Customers) onAllDeleted(event events.Event) {
	err := ths.DeleteAll()
	if err != nil {
		ths.log.Error("Can't delete customers: %v", err)
		return
	}
}

func NewCustomers(log logger.Logger, customersRepo repository.Customers, e *events.Events) (*Customers, error) {

	c := Customers{
		log:       log,
		customers: []models.CustomerData{},
		cache:     cache.New(10*time.Second, time.Minute),
		events:    e,
	}

	if e != nil {
		err := e.SubscribeOnBroadcast(repository.CustomerCreatedEvent, c.onCreated)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onCreated: %w", err)
		}

		err = e.SubscribeOnBroadcast(repository.CustomerAllDeletedEvent, c.onAllDeleted)
		if err != nil {
			return nil, fmt.Errorf("can't subscribe onAllDeleted: %w", err)
		}
	}

	persistentCustomers, err := customersRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get customers for cache: %w", err)
	}
	for _, pc := range persistentCustomers {
		err = c.Create(pc)
		if err != nil {
			return nil, fmt.Errorf("can't init customers cache: %w", err)
		}
	}

	return &c, nil
}
