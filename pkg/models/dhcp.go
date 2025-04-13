package models

import (
	"github.com/google/uuid"
	"time"
)

const (
	DhcpLeaseStatusAllocated = 0
	DhcpLeaseStatusAccepted  = 1

	DhcpPoolTypeNone   = 0
	DhcpPoolTypeShared = 10
	DhcpPoolTypeStatic = 20
)

type DhcpServerData struct {
	CommonData
	Name     string        `json:"name"`
	ServerId uuid.NullUUID `gorm:"column:server_id;type:uuid"`
	Ip       string        `json:"ip"`
	Descr    string        `json:"descr"`
	Type     int           `json:"type"`
}

func (DhcpServerData) TableName() string {
	return "ncc_dhcp_server"
}

type LeaseData struct {
	CommonData
	ServerId         uuid.NullUUID `gorm:"column:server_id;type:uuid"`
	CustomerId       uuid.NullUUID `gorm:"column:customer_id;type:uuid"`
	Customer         CustomerData  `gorm:"foreignKey:CustomerId"`
	Ip               string        `gorm:"column:ip;index"`
	Subnet           string        `json:"subnet"`
	Router           string        `json:"router"`
	Dns1             string        `json:"dns1"`
	Dns2             string        `json:"dns2"`
	Mac              string        `json:"mac"`
	Cvid             int           `json:"cvid"`
	Port             int           `json:"port"`
	Remote           string        `json:"remote"`
	Circuit          string        `json:"circuit"`
	LeaseTime        int64         `gorm:"column:leasetime"`
	Start            time.Time     `json:"start"`
	Expire           time.Time     `json:"expire"`
	Hostname         string        `json:"hostname"`
	DeviceId         uuid.NullUUID `gorm:"column:device_id;type:uuid"`
	Device           DeviceData    `gorm:"foreignKey:DeviceId"`
	Status           int           `json:"status"`
	MarkedForRemoval bool          `json:"marked_for_removal"`
	IfName           string        `gorm:"column:ifname"`
	LinkIndex        int           `gorm:"column:link_index"`
}

func (LeaseData) TableName() string {
	return "ncc_dhcp_lease"
}

type LeaseLogData struct {
	CommonData
	ServerId   uuid.NullUUID `gorm:"column:server_id;type:uuid"`
	CustomerId uuid.NullUUID `gorm:"column:customer_id;type:uuid"`
	Customer   CustomerData  `gorm:"foreignKey:CustomerId"`
	Ip         string        `gorm:"column:ip;index"`
	Subnet     string        `json:"subnet"`
	Router     string        `json:"router"`
	Dns1       string        `json:"dns1"`
	Dns2       string        `json:"dns2"`
	Mac        string        `json:"mac"`
	Cvid       int           `json:"cvid"`
	Port       int           `json:"port"`
	Remote     string        `json:"remote"`
	Circuit    string        `json:"circuit"`
	LeaseTime  int64         `gorm:"column:leasetime"`
	Start      time.Time     `json:"start"`
	Stop       time.Time     `json:"stop"`
	Expire     time.Time     `json:"expire"`
	Hostname   string        `json:"hostname"`
	DeviceId   uuid.NullUUID `gorm:"column:device_id;type:uuid"`
	Device     DeviceData    `gorm:"foreignKey:DeviceId"`
	IfName     string        `gorm:"column:ifname"`
	LinkIndex  int           `gorm:"column:link_index"`
}

func (LeaseLogData) TableName() string {
	return "ncc_dhcp_lease_log"
}

type DhcpPoolData struct {
	CommonData
	Name              string `json:"name"`
	Type              int    `json:"type"`
	RangeStart        string `json:"range_start"`
	RangeEnd          string `json:"range_end"`
	Mask              string `json:"mask"`
	Gateway           string `json:"gateway"`
	Dns1              string `json:"dns1"`
	Dns2              string `json:"dns2"`
	LeaseTime         int    `gorm:"column:leasetime"`
	UserChangeAllowed bool   `json:"user_change_allowed"`
}

func (DhcpPoolData) TableName() string {
	return "ncc_dhcp_pool"
}

type DhcpBindingData struct {
	CommonData
	CustomerId uuid.NullUUID `gorm:"column:customer_id;type:uuid;not null"`
	Customer   CustomerData  `gorm:"foreignKey:CustomerId"`
	DeviceId   uuid.NullUUID `gorm:"column:device_id;type:uuid"`
	Device     DeviceData    `gorm:"foreignKey:DeviceId"`
	Remote     string        `json:"remote"`
	Port       int           `json:"port"`
	Ip         string        `json:"ip"`
	Mac        string        `json:"mac"`
	Cvid       int           `json:"cvid"`
	PoolId     uuid.NullUUID `gorm:"column:pool_id;type:uuid"`
	Pool       DhcpPoolData  `gorm:"foreignKey:PoolId"`
}

func (DhcpBindingData) TableName() string {
	return "ncc_dhcp_binding"
}
