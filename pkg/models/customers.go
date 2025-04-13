package models

import (
	"database/sql"
	"github.com/google/uuid"
	"time"
)

const (
	CustomerStateActive  = 10
	CustomerStateBlocked = 20

	ServiceTypeInternet = "Доступ в Internet"
	ServiceTypeIptv     = "IPTV"
	ServiceTypeCatv     = "Кабельное ТВ"

	ContractTypePersonal   = 10
	ContractTypeEnterprise = 20
)

var (
	ServiceState = map[int]string{
		10: "Включен",
		20: "Выключено",
	}
)

type CustomerData struct {
	CommonData
	Uid      string  `gorm:"column:uid;unique;not null"`
	Login    string  `gorm:"column:login;unique;not null"`
	Phone    string  `gorm:"column:phone"`
	Deposit  float64 `gorm:"column:deposit;not null"`
	Credit   float64 `gorm:"column:credit;not null"`
	Password string  `gorm:"column:password"`
	Pin      string  `gorm:"column:pin"`
	Scores   int     `gorm:"column:scores"`
	Name     string  `gorm:"column:name"`

	FirstName  string `gorm:"column:first_name"`
	LastName   string `gorm:"column:last_name"`
	MiddleName string `gorm:"column:middle_name"`

	Build             string       `gorm:"column:build"`
	Flat              string       `gorm:"column:flat"`
	Email             string       `gorm:"column:email"`
	Comments          string       `gorm:"column:comments"`
	MonitoringEnabled bool         `gorm:"column:monitoring_enabled"`
	BlockingDate      time.Time    `gorm:"column:blocking_date"`
	BlockingTill      time.Time    `gorm:"column:blocking_till"`
	BlockingLastSet   time.Time    `gorm:"column:blocking_last_set"`
	BlockingState     int          `gorm:"column:blocking_state"`
	CreditExpire      sql.NullTime `gorm:"column:credit_expire"`
	CreditDaysLeft    int          `gorm:"column:credit_days_left"`

	GroupId uuid.NullUUID     `gorm:"column:group_id;type:uuid"`
	Group   CustomerGroupData `gorm:"foreignKey:GroupId"`

	CityId uuid.NullUUID `gorm:"type:uuid;column:city_id"`
	City   CityData      `gorm:"foreignKey:CityId"`

	StreetId uuid.NullUUID `gorm:"column:street_id;type:uuid"`
	Street   StreetData    `gorm:"foreignKey:StreetId"`

	ContractId uuid.NullUUID `gorm:"column:contract_id;type:uuid"`
	Contract   ContractData  `gorm:"foreignKey:ContractId"`

	ServiceInternetId    uuid.NullUUID       `gorm:"column:service_internet_id;type:uuid"`
	ServiceInternet      ServiceInternetData `gorm:"foreignKey:ServiceInternetId"`
	ServiceInternetState int                 `gorm:"column:service_internet_state"`

	ServiceIptvId    uuid.NullUUID   `gorm:"column:service_iptv_id;type:uuid"`
	ServiceIptv      ServiceIptvData `gorm:"foreignKey:ServiceIptvId"`
	ServiceIptvState int             `gorm:"column:service_iptv_state"`

	ServiceCatvId    uuid.NullUUID   `gorm:"column:service_catv_id;type:uuid"`
	ServiceCatv      ServiceCatvData `gorm:"foreignKey:ServiceCatvId"`
	ServiceCatvState int             `gorm:"column:service_catv_state"`

	Flags                     []CustomerFlagData          `gorm:"foreignKey:CustomerID;references:Id"`
	ServiceInternetCustomData []ServiceInternetCustomData `gorm:"foreignKey:CustomerId;references:Id"`

	VerifiedTs sql.NullTime `gorm:"column:verified_ts"`
	VerifiedBy string       `gorm:"column:verified_by"`

	BankAccounts []BankAccountData `gorm:"foreignKey:CustomerId"`
}

func (CustomerData) TableName() string {
	return "ncc_customer"
}

type CustomerGroupData struct {
	CommonData
	Name string `gorm:"column:name"`
}

func (CustomerGroupData) TableName() string {
	return "ncc_customer_group"
}

type ContractData struct {
	CommonData
	Type  int       `gorm:"column:type"`
	Date  time.Time `gorm:"column:date"`
	Phone string    `gorm:"column:phone"`
	Email string    `gorm:"column:email"`
	Name  string    `gorm:"column:name"`

	FirstName  string `gorm:"column:first_name"`
	LastName   string `gorm:"column:last_name"`
	MiddleName string `gorm:"column:middle_name"`

	Code                string    `gorm:"column:code"`
	Document            string    `gorm:"column:document"`
	DocumentDate        time.Time `gorm:"column:document_date"`
	Number              string    `gorm:"column:number"`
	Address             string    `gorm:"column:address"`
	Contact             string    `gorm:"column:contact"`
	DocumentSerial      string    `gorm:"column:document_serial"`
	DocumentNumber      string    `gorm:"column:document_number"`
	DocumentIssuedBy    string    `gorm:"column:document_issued_by"`
	BirthDate           time.Time `gorm:"column:birth_date"`
	KPP                 string    `gorm:"column:kpp"`
	INN                 string    `gorm:"column:inn"`
	OGRN                string    `gorm:"column:ogrn"`
	OGRNIP              string    `gorm:"column:ogrnip"`
	RegistrationAddress string    `gorm:"column:registration_address"`
	ResidentialAddress  string    `gorm:"column:residential_address"`
}

func (ContractData) TableName() string {
	return "ncc_contract"
}

type CustomerFlagData struct {
	CommonData
	CustomerID uuid.NullUUID `gorm:"column:customer_id;type:uuid;not null"`
	Customer   CustomerData  `gorm:"foreignKey:CustomerID"`
	Name       string        `gorm:"column:name"`
	Val        string        `gorm:"column:val"`
}

func (CustomerFlagData) TableName() string {
	return "ncc_customer_flag"
}

func StringState(val int) string {
	switch val {
	case CustomerStateActive:
		return "Active"
	case CustomerStateBlocked:
		return "Blocked"
	}
	return ""
}

type CustomerContact struct {
	CommonData
	CustomerId string       `gorm:"type:uuid;column:customer_id;not null"`
	Customer   CustomerData `gorm:"foreignKey:CustomerId"`
	CType      string       `gorm:"column:ctype;not null"`
	Sms        bool         `gorm:"column:sms"`
	Payments   bool         `gorm:"column:payments"`
	Promo      bool         `gorm:"column:promo"`
	Problems   bool         `gorm:"column:problems"`
	Data       string       `gorm:"column:data"`
}

func (CustomerContact) TableName() string {
	return "ncc_customer_contacts"
}

type CustomerServiceLink struct {
	CustomerId string `gorm:"type:uuid;column:customer_id;not null;index:idx_customer_service,unique"`
	ServiceId  string `gorm:"type:uuid;column:service_id;not null;index:idx_customer_service,unique"`
}

func (CustomerServiceLink) TableName() string {
	return "ncc_customer_service_link"
}

type DocumentTypeData struct {
	CommonData
	Name string `gorm:"column:name;unique"`
	Code string `gorm:"column:code;unique"`
}

func (DocumentTypeData) TableName() string {
	return "ncc_document_type"
}

type BankAccount struct {
	CommonData
	CustomerId string       `gorm:"type:uuid;column:customer_id;not null"`
	Customer   CustomerData `gorm:"foreignKey:CustomerId"`
	BIK        string       `gorm:"column:bik"`
	Number     string       `gorm:"column:number"`
	KPP        string       `gorm:"column:kpp"`
	BankName   string       `gorm:"column:bank_name"`
	CorrNumber string       `gorm:"column:corr_number"`
	Descr      string       `gorm:"column:descr"`
}

func (BankAccount) TableName() string {
	return "ncc_bank_account"
}
