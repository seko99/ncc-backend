package models

import "github.com/google/uuid"

type ScoreExchangeLogData struct {
	CommonData
	CustomerId     uuid.NullUUID    `gorm:"column:customer_id;type:uuid;not null"`
	Customer       CustomerData     `gorm:"foreignKey:CustomerId"`
	Amount         int              `gorm:"column:amount"`
	ScoreProductId uuid.NullUUID    `gorm:"column:score_product_id;type:uuid;not null"`
	ScoreProduct   ScoreProductData `gorm:"foreignKey:ScoreProductId"`
}

func (ScoreExchangeLogData) TableName() string {
	return "ncc_score_exchange_log"
}

type ScoreProductData struct {
	CommonData
	PaymentTypeId uuid.NullUUID     `gorm:"column:payment_type_id;type:uuid"`
	PaymentType   ScorePaymentTypes `gorm:"foreignKey:PaymentTypeId"`
	PaymentAmount float64           `gorm:"column:payment_amount"`
	Name          string            `gorm:"column:name"`
	Scores        int               `gorm:"column:scores"`
	Code          string            `gorm:"column:code"`
	Url           string            `gorm:"column:url"`
	Days          int               `gorm:"column:days"`
}

func (ScoreProductData) TableName() string {
	return "ncc_score_product"
}
