package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
)

type ServiceInternet struct {
	storage *psqlstorage.Storage
	events  *events.Events
}

func (s *ServiceInternet) GetById(id string) (*models2.ServiceInternetData, error) {
	var service *models2.ServiceInternetData

	r := s.storage.GetDB().Model(&models2.ServiceInternetData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = :id", sql.Named("id", id)).
		Find(&service)

	if r.Error != nil {
		return nil, r.Error
	}

	if r.RowsAffected > 0 {
		return service, nil
	}

	return nil, nil
}

func (s *ServiceInternet) Create(service models2.ServiceInternetData) error {
	r := s.storage.GetDB().Model(&models2.ServiceInternetData{}).
		Create(&service)

	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *ServiceInternet) Upsert(service models2.ServiceInternetData) error {
	r := s.storage.GetDB().Model(&models2.ServiceInternetData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&service)

	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *ServiceInternet) Get() ([]models2.ServiceInternetData, error) {
	var services []models2.ServiceInternetData

	r := s.storage.GetDB().Model(&models2.ServiceInternetData{}).
		Preload(clause.Associations).
		//Preload("ServiceInternetCustomData").
		Where("delete_ts is null").
		Find(&services)

	if r.Error != nil {
		return nil, r.Error
	}

	return services, nil
}

func (s *ServiceInternet) GetCustomDataMap() (map[string]models2.ServiceInternetCustomData, error) {
	customDataMap := map[string]models2.ServiceInternetCustomData{}

	var data []models2.ServiceInternetCustomData

	r := s.storage.GetDB().Model(&models2.ServiceInternetCustomData{}).
		Where("delete_ts is null").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	for _, d := range data {
		if !d.DeleteTs.IsZero() {
			continue
		}
		customDataMap[d.CustomerId.UUID.String()] = d
	}

	return customDataMap, nil
}

func (s *ServiceInternet) GetCustomDataByCustomer(customer models2.CustomerData) (*models2.ServiceInternetCustomData, error) {
	var customData models2.ServiceInternetCustomData

	r := s.storage.GetDB().Model(&models2.ServiceInternetCustomData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("customer_id = @customer_id", sql.Named("customer_id", customer.Id)).
		Find(&customData)
	if r.Error != nil {
		return nil, r.Error
	}

	return &customData, nil
}

func NewServiceInternet(storage *psqlstorage.Storage, e *events.Events) *ServiceInternet {
	services := &ServiceInternet{
		storage: storage,
		events:  e,
	}
	return services
}
