package models

import (
	"github.com/google/uuid"
	"time"
)

const (
	AccountTypeCash     = 10
	AccountTypeCashless = 20
)

type PaymentData struct {
	CommonData
	Pid           int       `json:"pid"`
	Date          time.Time `json:"date"`
	Amount        float64   `json:"amount"`
	Descr         string    `json:"descr"`
	DepositBefore float64   `json:"deposit_before"`

	PaymentTypeId uuid.NullUUID   `gorm:"column:payment_type_id;type:uuid;not null"`
	PaymentType   PaymentTypeData `gorm:"foreignKey:PaymentTypeId"`

	CustomerId uuid.NullUUID `gorm:"column:customer_id;type:uuid;not null"`
	Customer   CustomerData  `gorm:"foreignKey:CustomerId"`
}

func (PaymentData) TableName() string {
	return "ncc_payment"
}

type PaymentTypeData struct {
	CommonData
	Name          string `json:"name"`
	AccountType   int    `json:"account_type"`
	ManualEnabled bool   `json:"manual_enabled"`
	Code          string `json:"code"`
}

func (PaymentTypeData) TableName() string {
	return "ncc_payment_type"
}

type PaymentSystemData struct {
	CommonData
	Enabled  bool   `json:"enabled"`
	Name     string `json:"name"`
	Token    string `json:"token"`
	TestMode bool   `json:"test_mode"`

	PaymentTypeId uuid.NullUUID   `gorm:"column:payment_type_id;type:uuid;not null"`
	PaymentType   PaymentTypeData `gorm:"foreignKey:PaymentTypeId"`

	TestCustomerId uuid.NullUUID `gorm:"column:test_customer_id;type:uuid"`
	TestCustomer   CustomerData  `gorm:"foreignKey:TestCustomerId"`

	UserId uuid.NullUUID `gorm:"column:user_id;type:uuid"`
	User   SecUserData   `gorm:"foreignKey:UserId"`
}

func (PaymentSystemData) TableName() string {
	return "ncc_payment_system"
}
