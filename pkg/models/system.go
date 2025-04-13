package models

import "time"

type SecUserData struct {
	CommonData
	Login                 string `json:"login"`
	LoginLc               string `gorm:"column:login_lc"`
	Password              string `json:"password"`
	PasswordEncryption    string `json:"password_encryption"`
	Name                  string `json:"name"`
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	MiddleName            string `json:"middle_name"`
	Position              string `gorm:"column:position_"`
	Email                 string `json:"email"`
	Language              string `gorm:"column:language_"`
	TimeZone              string `json:"time_zone"`
	TimeZoneAuto          bool   `json:"time_zone_auto"`
	Active                bool   `json:"active"`
	GroupId               string `gorm:"type:uuid;column:group_id"`
	IPMask                string `gorm:"column:ip_mask"`
	ChangePasswordAtLogon bool   `json:"change_password_at_logon"`
}

func (SecUserData) TableName() string {
	return "sec_user"
}

type SecGroupData struct {
	CommonData
	Name     string `json:"name"`
	ParentId string `gorm:"type:uuid;column:parent_id"`
}

func (SecGroupData) TableName() string {
	return "sec_group"
}

type SysServer struct {
	CommonData
	Name      string `gorm:"type:varchar(255);unique;index"`
	IsRunning bool
	Data      string `gorm:"type:text"`
}

func (SysServer) TableName() string {
	return "sys_server"
}

type SysConfig struct {
	CommonData
	Name  string `gorm:"type:varchar(255);unique;index"`
	Value string `gorm:"column:value_;type:text;not null"`
}

func (SysConfig) TableName() string {
	return "sys_config"
}

type SysAccessToken struct {
	CommonData
	TokenValue          string
	TokenBytes          []byte
	AuthenticationKey   string
	AuthenticationBytes []byte
	Expiry              time.Time
	UserLogin           string
	Locale              string
	RefreshTokenValue   string
}

func (SysAccessToken) TableName() string {
	return "sys_access_token"
}
