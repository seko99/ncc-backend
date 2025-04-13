package models

import "github.com/google/uuid"

type SormCustomersErrorsData struct {
	CommonData

	Login      string        `gorm:"column:login"`
	CustomerId uuid.NullUUID `gorm:"column:customer_id;type:uuid;not null"`
	Customer   CustomerData  `gorm:"foreignKey:CustomerId"`
	Reason     string        `gorm:"column:reason"`
}

func (SormCustomersErrorsData) TableName() string {
	return "ncc_sorm_customers_errors"
}
