package models

import "github.com/google/uuid"

type ScorePaymentTypes struct {
	CommonData
	PaymentTypeId     uuid.NullUUID   `gorm:"column:payment_type_id;type:uuid;not null"`
	PaymentType       PaymentTypeData `gorm:"foreignKey:PaymentTypeId"`
	PaymentBaseAmount float64         `json:"payment_base_amount"`
	ScoresAmount      int             `json:"scores_amount"`
}

func (ScorePaymentTypes) TableName() string {
	return "ncc_score_payment_types"
}
