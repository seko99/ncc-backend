package models

import (
	"github.com/google/uuid"
	"time"
)

type SessionData struct {
	CommonData
	AcctSessionId string    `gorm:"column:acct_session_id;unique;not null;index"`
	Login         string    `gorm:"column:login;not null;index"`
	StartTime     time.Time `json:"start_time"`
	Duration      int64     `json:"duration"`
	NasName       string    `json:"nas_name"`
	Ip            string    `gorm:"column:ip;unique;not null;index"`
	Mac           string    `gorm:"column:mac"`
	Circuit       string    `json:"circuit"`
	Remote        string    `json:"remote"`
	ServiceName   string    `json:"service_name"`
	LastAlive     time.Time `json:"last_alive"`
	OctetsIn      int64     `json:"octets_in"`
	OctetsOut     int64     `json:"octets_out"`

	CustomerId uuid.NullUUID `gorm:"column:customer_id;type:uuid;not null"`
	Customer   CustomerData  `gorm:"foreignKey:CustomerId"`

	NasId   uuid.NullUUID `gorm:"column:nas_id;type:uuid"`
	Nas     NasData       `gorm:"foreignKey:NasId"`
	NasPort uint32        `json:"nas_port"`

	ServiceInternetId uuid.NullUUID       `gorm:"column:service_internet_id;type:uuid;not null"`
	ServiceInternet   ServiceInternetData `gorm:"foreignKey:ServiceInternetId"`
}

func (SessionData) TableName() string {
	return "ncc_session"
}

func (SessionData) FromSessionLog(sessionLog SessionsLogData) SessionData {
	return SessionData{
		CommonData:        sessionLog.CommonData,
		AcctSessionId:     sessionLog.AcctSessionId,
		Login:             sessionLog.Login,
		StartTime:         sessionLog.StartTime,
		Duration:          sessionLog.Duration,
		NasName:           sessionLog.NasName,
		Ip:                sessionLog.Ip,
		Mac:               sessionLog.Mac,
		Circuit:           sessionLog.Circuit,
		Remote:            sessionLog.Remote,
		ServiceName:       sessionLog.ServiceName,
		LastAlive:         sessionLog.LastAlive,
		OctetsIn:          sessionLog.OctetsIn,
		OctetsOut:         sessionLog.OctetsOut,
		CustomerId:        sessionLog.CustomerId,
		NasId:             sessionLog.NasId,
		ServiceInternetId: sessionLog.ServiceInternetId,
	}
}
