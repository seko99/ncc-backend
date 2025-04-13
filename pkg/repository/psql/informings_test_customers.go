package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"gorm.io/gorm/clause"
)

type InformingsTestCustomers struct {
	storage *psqlstorage.Storage
}

func (s *InformingsTestCustomers) Get() ([]models.InformingTestCustomerData, error) {
	var customers []models.InformingTestCustomerData
	r := s.storage.GetDB().Model(models.InformingTestCustomerData{}).
		Preload(clause.Associations).
		Preload("Customer.Group").
		Preload("Customer.ServiceInternet").
		Where("delete_ts is null").
		Find(&customers)
	if r.Error != nil {
		return nil, r.Error
	}
	return customers, nil
}

func NewInformingsTestCustomers(storage *psqlstorage.Storage) *InformingsTestCustomers {
	return &InformingsTestCustomers{
		storage: storage,
	}
}
