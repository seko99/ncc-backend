package models

import (
	"github.com/google/uuid"
	"time"
)

type AccountData struct {
	CommonData
	ParentId      uuid.NullUUID `gorm:"type:uuid;column:parent_id"`
	Name          string
	Type          int
	Code          string
	Remain        float64
	LastOperation time.Time
	Descr         string
}

func (AccountData) TableName() string {
	return "ncc_account"
}

type AccountTransactionData struct {
	CommonData
	Tid             uint64
	SrcAccountId    uuid.NullUUID `gorm:"type:uuid;column:src_account_id;not null"`
	SrcAccount      AccountData   `gorm:"foreignKey:SrcAccountId"`
	DstAccountId    uuid.NullUUID `gorm:"type:uuid;column:dst_account_id;not null"`
	DstAccount      AccountData   `gorm:"foreignKey:DstAccountId"`
	Amount          float64
	SrcRemainBefore float64
	DstRemainBefore float64
	Descr           string
	PaymentId       uuid.NullUUID `gorm:"type:uuid;column:payment_id"`
	Payment         PaymentData   `gorm:"foreignKey:PaymentId"`
}

func (AccountTransactionData) TableName() string {
	return "ncc_account_transaction"
}

type BankAccountData struct {
	CommonData
	CustomerId uuid.NullUUID `gorm:"type:uuid;column:customer_id;not null"`
	Customer   CustomerData  `gorm:"foreignKey:CustomerId"`
	BIK        string
	Number     string
	KPP        string
	BankName   string
	CorrNumber string
	Descr      string
}

func (BankAccountData) TableName() string {
	return "ncc_bank_account"
}
