package models

import "github.com/google/uuid"

type ScoreLogData struct {
	CommonData
	CustomerId   uuid.NullUUID `gorm:"column:customer_id;type:uuid;not null"`
	Customer     CustomerData  `gorm:"foreignKey:CustomerId"`
	Scores       int           `json:"scores"`
	PaymentId    uuid.NullUUID `gorm:"column:payment_id;type:uuid;not null"`
	Payment      PaymentData   `gorm:"foreignKey:PaymentId"`
	ScoresBefore int           `json:"scores_before"`
}

func (ScoreLogData) TableName() string {
	return "ncc_score_log"
}
