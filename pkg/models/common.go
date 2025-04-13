package models

import (
	"database/sql"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

func NewNullUUID(s string) uuid.NullUUID {
	if len(s) == 0 {
		return uuid.NullUUID{}
	}
	var n uuid.NullUUID
	_ = n.Scan(s)
	return n
}

type CommonData struct {
	Id        string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Version   int       `json:"version"`
	CreateTs  time.Time `json:"create_ts"`
	CreatedBy string    `json:"created_by"`
	UpdateTs  time.Time `json:"update_ts"`
	UpdatedBy string    `json:"updated_by"`
	DeleteTs  time.Time `gorm:"default:null"`
	DeletedBy string    `gorm:"default:null"`
}

func (c *CommonData) BeforeCreate(tx *gorm.DB) error {
	if len(c.Id) == 0 {
		c.Id = uuid.NewString()
	}
	if c.CreateTs.IsZero() {
		c.CreateTs = time.Now()
	}
	return nil
}

func (c *CommonData) BeforeUpdate(tx *gorm.DB) error {
	if c.UpdateTs.IsZero() {
		c.UpdateTs = time.Now()
	}
	return nil
}
