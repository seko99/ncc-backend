package models

import (
	"github.com/google/uuid"
	"time"
)

const (
	InformingModeRegular = 10
	InformingModeTest    = 20

	InformingStateEnabled  = 10
	InformingStateDisabled = 20

	InformingRepeatingDaily   = 10
	InformingRepeatingMonthly = 20
	InformingRepeatingNever   = 0

	FieldDeposit       = "deposit"
	FieldCredit        = "credit"
	FieldLogin         = "login"
	FieldGroup         = "group"
	FieldVerified      = "verified"
	FieldSent          = "sentFlag"
	FieldInternetState = "internetState"
	FieldBlockingState = "blockingState"

	ExprFalse = "false"
	ExprTrue  = "true"

	ExprEq = "="
	ExprNe = "!="
	ExprLt = "<"
	ExprLe = "<="
	ExprGt = ">"
	ExprGe = ">="
)

type InformingData struct {
	CommonData
	Message    string                   `gorm:"column:message"`
	Start      time.Time                `gorm:"column:start"`
	Type       string                   `gorm:"column:itype"`
	Condition  string                   `gorm:"column:condition"`
	Conditions []InformingConditionData `gorm:"foreignKey:InformingID"`
	Descr      string                   `gorm:"column:descr"`
	Name       string                   `gorm:"column:name"`
	State      int                      `gorm:"column:state"`
	Repeating  int                      `gorm:"column:repeating"`
	Mode       int                      `gorm:"column:mode"`
}

func (InformingData) TableName() string {
	return "ncc_informing"
}

type InformingConditionData struct {
	CommonData
	InformingID uuid.NullUUID `gorm:"column:informing_id;type:uuid;not null"`
	Informing   InformingData `gorm:"foreignKey:InformingID"`
	Field       string        `json:"field"`
	Expr        string        `json:"expr"`
	Val         string        `json:"val"`
}

func (InformingConditionData) TableName() string {
	return "ncc_informing_condition"
}

type InformingTestCustomerData struct {
	CommonData
	CustomerId uuid.NullUUID `gorm:"column:customer_id;type:uuid;not null"`
	Customer   CustomerData  `gorm:"foreignKey:CustomerId"`
}

func (InformingTestCustomerData) TableName() string {
	return "ncc_informing_test_customers"
}

type InformingLogData struct {
	CommonData
	InformingId string        `gorm:"column:informing_id;type:uuid;not null"`
	Informing   InformingData `gorm:"foreignKey:InformingId"`
	CustomerId  string        `gorm:"column:customer_id;type:uuid;not null"`
	Customer    CustomerData  `gorm:"foreignKey:CustomerId"`
	Phone       string        `gorm:"column:phone;not null"`
	Message     string        `gorm:"column:message;not null"`
	Status      int           `gorm:"column:status"`
}

func (InformingLogData) TableName() string {
	return "ncc_informing_log"
}
