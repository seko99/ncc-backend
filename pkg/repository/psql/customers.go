package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"fmt"
	"gorm.io/gorm/clause"
	"time"
)

type Customers struct {
	storage *psqlstorage.Storage
	events  *events.Events
}

func (s *Customers) IsBlocked(state int) bool {
	switch state {
	case 20:
		return true
	}
	return false
}

func (s *Customers) Get(limit ...int) ([]models.CustomerData, error) {
	var customers []models.CustomerData

	db := s.storage.GetDB().Model(models.CustomerData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Order("login")

	if len(limit) > 0 {
		db.Limit(limit[0])
	}

	r := db.Find(&customers)

	if r.Error != nil {
		return nil, r.Error
	}

	return customers, nil
}

func (s *Customers) GetById(id string) (*models.CustomerData, error) {
	var customer *models.CustomerData

	r := s.storage.GetDB().Model(models.CustomerData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Find(&customer)

	if r.Error != nil {
		return nil, r.Error
	}

	return customer, nil
}

func (s *Customers) GetByUid(uid string) (*models.CustomerData, error) {
	var customer *models.CustomerData

	r := s.storage.GetDB().Model(models.CustomerData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("uid = @uid", sql.Named("uid", uid)).
		Find(&customer)

	if r.Error != nil {
		return nil, r.Error
	}

	return customer, nil
}

func (s *Customers) GetByLogin(login string) (*models.CustomerData, error) {
	var customer *models.CustomerData

	r := s.storage.GetDB().Model(models.CustomerData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("login = @login", sql.Named("login", login)).
		Find(&customer)

	if r.Error != nil {
		return nil, r.Error
	}

	return customer, nil
}

func (s *Customers) GetByFeeAmountAndSessions(amount float64, sdate string) ([]models.CustomerData, error) {
	var customer []models.CustomerData

	r := s.storage.GetDB().Raw(fmt.Sprintf(`select c.*
	from ncc_customer c
    left join ncc_session_log nsl on c.id = nsl.customer_id
    left join ncc_customer_fee_log ncfl on c.id = ncfl.customer_id
where
    ncfl.fee_amount=%0.2f
    and ncfl.fee_timestamp>'%s'
    and nsl.stop_time>'%s'
group by c.id, c.login, c.uid`, amount, sdate, sdate)).Scan(&customer)

	/*	r := s.storage.GetDB().Model(models.CustomerData{}).
		Joins("LEFT JOIN ? ON ?.? = ?.?",
			clause.Table{Name: "ncc_session_log"},
			clause.Table{Name: "ncc_session_log"}, clause.Column{Name: "customer_id"},
			clause.Table{Name: "ncc_customer"}, clause.Column{Name: "id"}).
		Joins("LEFT JOIN ? ON ?.? = ?.?",
			clause.Table{Name: "ncc_customer_fee_log"},
			clause.Table{Name: "ncc_customer_fee_log"}, clause.Column{Name: "customer_id"},
			clause.Table{Name: "ncc_customer"}, clause.Column{Name: "id"}).
		Where("?.fee_amount = ?",
			clause.Table{Name: "ncc_customer_fee_log"}, amount).
		Where("?.fee_timestamp > ?",
			clause.Table{Name: "ncc_customer_fee_log"}, sdate).
		Where("?.stop_time > ?",
			clause.Table{Name: "ncc_session_log"}, sdate).
		Find(&customer)
	*/
	if r.Error != nil {
		return nil, r.Error
	}

	return customer, nil
}

func (s *Customers) GetByState(state int) ([]models.CustomerData, error) {
	var customers []models.CustomerData

	r := s.storage.GetDB().Model(models.CustomerData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("blocking_state = @state", sql.Named("state", state)).
		Find(&customers)

	if r.Error != nil {
		return nil, r.Error
	}

	return customers, nil
}

func (s *Customers) GetByGroup(id string) ([]models.CustomerData, error) {
	var customers []models.CustomerData

	r := s.storage.GetDB().Model(models.CustomerData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("group_id = @id", sql.Named("id", id)).
		Find(&customers)

	if r.Error != nil {
		return nil, r.Error
	}

	return customers, nil
}

func (s *Customers) GetGroupByName(name string) (*models.CustomerGroupData, error) {
	var group *models.CustomerGroupData

	r := s.storage.GetDB().Model(models.CustomerGroupData{}).
		Where("delete_ts is null").
		Where("name = @name", sql.Named("name", name)).
		Find(&group)

	if r.Error != nil {
		return nil, r.Error
	}

	return group, nil
}

func (s *Customers) SetDeposit(id string, deposit float64) error {

	r := s.storage.GetDB().Model(&models.CustomerData{}).
		Where("id = @id", sql.Named("id", id)).
		Update("deposit", deposit)

	if r.Error != nil {
		return r.Error
	}

	if s.events != nil {
		err := s.events.PublishEvent(events.NewEvent(repository.CustomerUpdatedEvent, models.CustomerData{
			CommonData: models.CommonData{
				Id: id,
			},
			Deposit: deposit,
		}))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.CustomerUpdatedEvent, err)
		}
	}

	return nil
}

func (s *Customers) SetCredit(id string, credit float64) error {

	r := s.storage.GetDB().Model(&models.CustomerData{}).
		Where("id = @id", sql.Named("id", id)).
		Update("credit", credit)

	if r.Error != nil {
		return r.Error
	}

	if s.events != nil {
		err := s.events.PublishEvent(events.NewEvent(repository.CustomerUpdatedEvent, models.CustomerData{
			CommonData: models.CommonData{
				Id: id,
			},
			Credit: credit,
		}))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.CustomerUpdatedEvent, err)
		}
	}

	return nil
}

func (s *Customers) SetCreditExpire(id string, creditExpire time.Time) error {

	r := s.storage.GetDB().Model(&models.CustomerData{}).
		Where("id = @id", sql.Named("id", id)).
		Update("credit_expire", creditExpire)

	if r.Error != nil {
		return r.Error
	}

	if s.events != nil {
		err := s.events.PublishEvent(events.NewEvent(repository.CustomerUpdatedEvent, models.CustomerData{
			CommonData: models.CommonData{
				Id: id,
			},
			CreditExpire: models.NewNullTime(creditExpire),
		}))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.CustomerUpdatedEvent, err)
		}
	}

	return nil
}

func (s *Customers) SetCreditDaysLeft(id string, creditDaysLeft int) error {

	r := s.storage.GetDB().Model(&models.CustomerData{}).
		Where("id = @id", sql.Named("id", id)).
		Update("credit_days_left", creditDaysLeft)

	if r.Error != nil {
		return r.Error
	}

	if s.events != nil {
		err := s.events.PublishEvent(events.NewEvent(repository.CustomerUpdatedEvent, models.CustomerData{
			CommonData: models.CommonData{
				Id: id,
			},
			CreditDaysLeft: creditDaysLeft,
		}))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.CustomerUpdatedEvent, err)
		}
	}

	return nil
}

func (s *Customers) SetState(id string, state int) error {

	r := s.storage.GetDB().Model(&models.CustomerData{}).
		Where("id = @id", sql.Named("id", id)).
		Update("blocking_state", state)

	if r.Error != nil {
		return r.Error
	}

	if s.events != nil {
		err := s.events.PublishEvent(events.NewEvent(repository.CustomerUpdatedEvent, models.CustomerData{
			CommonData: models.CommonData{
				Id: id,
			},
			BlockingState: state,
		}))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.CustomerUpdatedEvent, err)
		}
	}

	return nil
}

func (s *Customers) SetServiceInternetState(id string, state int) error {

	r := s.storage.GetDB().Model(&models.CustomerData{}).
		Where("id = @id", sql.Named("id", id)).
		Update("service_internet_state", state)

	if r.Error != nil {
		return r.Error
	}

	if s.events != nil {
		err := s.events.PublishEvent(events.NewEvent(repository.CustomerUpdatedEvent, models.CustomerData{
			CommonData: models.CommonData{
				Id: id,
			},
			ServiceInternetState: state,
		}))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.CustomerUpdatedEvent, err)
		}
	}

	return nil
}

func (s *Customers) GetFlag(flag models.CustomerFlagData) (models.CustomerFlagData, error) {
	var data models.CustomerFlagData
	r := s.storage.GetDB().Model(&models.CustomerFlagData{}).
		Where("customer_id = ?", flag.CustomerID).
		Where("name = ?", flag.Name).
		First(&data)

	if r.Error != nil {
		return models.CustomerFlagData{}, r.Error
	}

	return data, nil
}

func (s *Customers) SetFlag(customer models.CustomerData, flag models.CustomerFlagData) error {
	/*	r := s.storage.GetDB().Model(&models.CustomerFlagData{}).
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "customer_id"}, {Name: "name"}},
				DoUpdates: clause.AssignmentColumns([]string{"val"}),
			}).
			Create(&flag)

		if r.Error != nil {
			return r.Error
		}
	*/

	currentFlag, err := s.GetFlag(flag)
	if err != nil {
		r := s.storage.GetDB().Model(&models.CustomerFlagData{}).Create(&flag)

		if r.Error != nil {
			return r.Error
		}

		return nil
	}

	if currentFlag.Val != flag.Val {
		r := s.storage.GetDB().Model(&models.CustomerFlagData{}).
			Where("customer_id = ?", flag.CustomerID).
			Where("name = ?", flag.Name).
			Update("val", flag.Val)

		if r.Error != nil {
			return r.Error
		}
	}

	return nil
}

func (s *Customers) Create(data models.CustomerData) error {
	r := s.storage.GetDB().Model(&models.CustomerData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	if s.events != nil {
		err := s.events.PublishEvent(events.NewEvent(repository.CustomerCreatedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.CustomerCreatedEvent, err)
		}
	}

	return nil
}

func (s *Customers) Update(data models.CustomerData) error {
	r := s.storage.GetDB().Model(&models.CustomerData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}

	if s.events != nil {
		err := s.events.PublishEvent(events.NewEvent(repository.CustomerUpdatedEvent, data))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.CustomerUpdatedEvent, err)
		}
	}

	return nil
}

func (s *Customers) DeleteAll() error {
	r := s.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_customer")
	if r.Error != nil {
		return fmt.Errorf("can't delete: %w", r.Error)
	}

	if s.events != nil {
		err := s.events.PublishEvent(events.NewEvent(repository.CustomerAllDeletedEvent, nil))
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", repository.CustomerAllDeletedEvent, err)
		}
	}

	return nil
}

func NewCustomers(storage *psqlstorage.Storage, e *events.Events) *Customers {
	customers := Customers{
		storage: storage,
		events:  e,
	}
	return &customers
}
