package models

import (
	"github.com/google/uuid"
	"time"
)

type SnapshotData struct {
	CommonData
	Uid             string    `json:"uid"`
	Login           string    `json:"login"`
	Deposit         float64   `json:"deposit"`
	Credit          float64   `json:"credit"`
	Scores          int       `json:"scores"`
	BlockingDate    time.Time `json:"blocking_date"`
	BlockingTill    time.Time `json:"blocking_till"`
	BlockingLastSet time.Time `json:"blocking_last_set"`
	BlockingState   int       `json:"blocking_state"`
	CreditExpire    time.Time `json:"credit_expire"`
	CreditDaysLeft  int       `json:"credit_days_left"`

	GroupId string `json:"group_id"`

	ServiceInternetId       uuid.NullUUID       `gorm:"column:service_internet_id;type:uuid"`
	ServiceInternet         ServiceInternetData `gorm:"foreignJKey:ServiceInternetId"`
	ServiceInternetState    int                 `json:"service_internet_state"`
	ServiceInternetFee      float64             `json:"service_internet_fee"`
	ServiceInternetIP       string              `json:"service_internet_ip"`
	ServiceInternetIPFee    float64             `json:"service_internet_ip_fee"`
	ServiceInternetSpeedIn  int                 `json:"service_internet_speed_in"`
	ServiceInternetSpeedOut int                 `json:"service_internet_speed_out"`

	ServiceIptvId    uuid.NullUUID   `gorm:"column:service_iptv_id;type:uuid"`
	ServiceIptv      ServiceIptvData `gorm:"foreignKey:ServiceIptvId"`
	ServiceIptvState int             `json:"service_iptv_state"`
	ServiceIptvFee   float64         `json:"service_iptv_fee"`

	ServiceCatvId    uuid.NullUUID   `gorm:"column:service_catv_id;type:uuid"`
	ServiceCatv      ServiceCatvData `gorm:"foreignKey:ServiceCatvId"`
	ServiceCatvState int             `json:"service_catv_state"`
	ServiceCatvFee   float64         `json:"service_catv_fee"`

	HasSessions bool `json:"has_sessions"`
}

func (SnapshotData) TableName() string {
	return "ncc_snapshots"
}
