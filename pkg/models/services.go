package models

import (
	"github.com/google/uuid"
	"time"
)

const (
	ServiceStateEnabled   = 10
	ServiceStateDisabled  = 20
	ServiceStateSuspended = 30
	ServiceStateBlocked   = 40

	FeeTypeDaily   = 10
	FeeTypeMonthly = 20

	Speed1G   = 1024000
	Speed100M = 102400
	Speed50M  = 51200
	Speed30M  = 30720
	Speed10M  = 10240
)

type CommonServiceData struct {
	Name                   string        `json:"name"`
	Fee                    float64       `json:"fee"`
	FeeMethod              int           `json:"fee_method"`
	FeeType                int           `json:"fee_type"`
	FeeDate                int           `json:"fee_date"`
	ArrangeTermsChange     bool          `json:"arrange_terms_change"`
	ArrangeServiceChange   bool          `json:"arrange_service_change"`
	TermsChangeDate        time.Time     `json:"terms_change_date"`
	TermsChangeFee         float64       `json:"terms_change_fee"`
	ServiceChangeDate      time.Time     `json:"service_change_date"`
	ServiceChangeServiceId uuid.NullUUID `gorm:"column:service_change_service_id;type:uuid"`
	NegativeDepositAllowed bool          `json:"negative_deposit_allowed"`
	Code                   string        `json:"code"`
}

type ServiceInternetData struct {
	CommonData
	CommonServiceData

	SpeedIn  int     `json:"speed_in"`
	SpeedOut int     `json:"speed_out"`
	Ip       string  `json:"ip"`
	IpFee    float64 `json:"ip_fee"`

	TermsChangeSpeedIn  int     `json:"terms_change_speed_in"`
	TermsChangeSpeedOut int     `json:"terms_change_speed_out"`
	TermsChangeIpFee    float64 `json:"terms_change_ip_fee"`

	IPPoolId uuid.NullUUID `gorm:"type:uuid;column:ip_pool_id"`
	IPPool   IpPoolData    `gorm:"foreignKey:IPPoolId"`
}

func (ServiceInternetData) TableName() string {
	return "ncc_service_internet"
}

type ServiceInternetCustomData struct {
	CommonData

	Fee      float64 `json:"fee"`
	SpeedIn  int     `json:"speed_in"`
	SpeedOut int     `json:"speed_out"`
	Ip       string  `json:"ip"`
	IpFee    float64 `json:"ip_fee"`
	Mac      string  `json:"mac"`

	CustomerId uuid.NullUUID `gorm:"column:customer_id;type:uuid;not null"`
	Customer   CustomerData  `gorm:"foreignKey:CustomerId"`

	ServiceInternetId uuid.NullUUID       `gorm:"type:uuid;column:service_internet_id"`
	ServiceInternet   ServiceInternetData `gorm:"foreignKey:ServiceInternetId"`

	IpPoolId uuid.NullUUID `gorm:"column:ip_pool_id;type:uuid"`
	IpPool   IpPoolData    `gorm:"foreignKey:IpPoolId"`

	ArrangeServiceChange bool      `json:"arrange_service_change"`
	ServiceChangeDate    time.Time `json:"service_change_date"`

	ServiceChangeServiceId uuid.NullUUID `gorm:"column:service_change_service_id;type:uuid"`
	//ServiceChangeService   ServiceInternetData `gorm:"foreignKey:ServiceChangeServiceId"`
}

func (ServiceInternetCustomData) TableName() string {
	return "ncc_service_internet_custom_data"
}

type ServiceIptvData struct {
	CommonData
	CommonServiceData
}

func (ServiceIptvData) TableName() string {
	return "ncc_service_iptv"
}

type ServiceCatvData struct {
	CommonData
	CommonServiceData
}

func (ServiceCatvData) TableName() string {
	return "ncc_service_catv"
}

type ServiceData struct {
	CommonData
	Code string
	Name string
}

func (ServiceData) TableName() string {
	return "ncc_service"
}
