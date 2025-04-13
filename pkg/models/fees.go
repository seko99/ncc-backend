package models

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type FeeLogData struct {
	CommonData

	CustomerId uuid.NullUUID `gorm:"column:customer_id;type:uuid;not null"`
	Customer   CustomerData  `gorm:"foreignKey:CustomerId"`

	FeeTimestamp time.Time `json:"fee_timestamp"`
	FeeAmount    float64   `json:"fee_amount"`

	PrevDeposit  float64      `json:"prev_deposit"`
	NewDeposit   float64      `json:"new_deposit"`
	Credit       float64      `json:"credit"`
	CreditExpire sql.NullTime `json:"credit_expire"`

	ServiceInternetId uuid.NullUUID       `gorm:"column:service_internet_id;type:uuid"`
	ServiceInternet   ServiceInternetData `gorm:"foreignKey:ServiceInternetId"`

	ServiceIptvId uuid.NullUUID   `gorm:"column:service_iptv_id;type:uuid"`
	ServiceIptv   ServiceIptvData `gorm:"foreignKey:ServiceIptvId"`

	ServiceCatvId uuid.NullUUID   `gorm:"column:service_catv_id;type:uuid"`
	ServiceCatv   ServiceCatvData `gorm:"foreignKey:ServiceCatvId"`

	Descr string `json:"descr"`
}

func (FeeLogData) TableName() string {
	return "ncc_customer_fee_log"
}
