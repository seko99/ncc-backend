package models

import "github.com/google/uuid"

type SormCustomerServicesErrorsData struct {
	CommonData

	Login      string        `gorm:"column:login"`
	CustomerId uuid.NullUUID `gorm:"column:customer_id;type:uuid;not null"`
	Customer   CustomerData  `gorm:"foreignKey:CustomerId"`
	Reason     string        `gorm:"column:reason"`
}

func (SormCustomerServicesErrorsData) TableName() string {
	return "ncc_sorm_customer_services_errors"
}
