package models

import (
	"github.com/google/uuid"
	"layeh.com/radius/rfc2866"
	"time"
)

const (
	TerminateCauseUserRequest        = "User-Request"
	TerminateCauseLostCarrier        = "Lost-Carrier"
	TerminateCauseLostService        = "Lost-Service"
	TerminateCauseIdleTimeout        = "Idle-Timeout"
	TerminateCauseSessionTimeout     = "Session-Timeout"
	TerminateCauseAdminReset         = "Admin-Reset"
	TerminateCauseAdminReboot        = "Admin-Reboot"
	TerminateCausePortError          = "Port-Error"
	TerminateCauseNASError           = "NAS-Error"
	TerminateCauseNASRequest         = "NAS-Request"
	TerminateCauseNASReboot          = "NAS-Reboot"
	TerminateCausePortUnneeded       = "Port-Unneeded"
	TerminateCausePortPreempted      = "Port-Preempted"
	TerminateCausePortSuspended      = "Port-Suspended"
	TerminateCauseServiceUnavailable = "Service-Unavailable"
	TerminateCauseCallback           = "Callback"
	TerminateCauseUserError          = "User-Error"
	TerminateCauseHostRequest        = "Host-Request"
)

type TerminateCause uint32

func (c TerminateCause) String() string {
	switch c {
	case TerminateCause(rfc2866.AcctTerminateCause_Value_UserRequest):
		return TerminateCauseUserRequest
	case TerminateCause(rfc2866.AcctTerminateCause_Value_LostCarrier):
		return TerminateCauseLostCarrier
	case TerminateCause(rfc2866.AcctTerminateCause_Value_LostService):
		return TerminateCauseLostService
	case TerminateCause(rfc2866.AcctTerminateCause_Value_UserError):
		return TerminateCauseUserError
	case TerminateCause(rfc2866.AcctTerminateCause_Value_IdleTimeout):
		return TerminateCauseIdleTimeout
	case TerminateCause(rfc2866.AcctTerminateCause_Value_SessionTimeout):
		return TerminateCauseSessionTimeout
	case TerminateCause(rfc2866.AcctTerminateCause_Value_AdminReset):
		return TerminateCauseAdminReset
	case TerminateCause(rfc2866.AcctTerminateCause_Value_AdminReboot):
		return TerminateCauseAdminReboot
	case TerminateCause(rfc2866.AcctTerminateCause_Value_PortError):
		return TerminateCausePortError
	case TerminateCause(rfc2866.AcctTerminateCause_Value_NASError):
		return TerminateCauseNASError
	case TerminateCause(rfc2866.AcctTerminateCause_Value_NASRequest):
		return TerminateCauseNASRequest
	case TerminateCause(rfc2866.AcctTerminateCause_Value_NASReboot):
		return TerminateCauseNASReboot
	case TerminateCause(rfc2866.AcctTerminateCause_Value_PortUnneeded):
		return TerminateCausePortUnneeded
	case TerminateCause(rfc2866.AcctTerminateCause_Value_PortPreempted):
		return TerminateCausePortPreempted
	case TerminateCause(rfc2866.AcctTerminateCause_Value_PortSuspended):
		return TerminateCausePortSuspended
	case TerminateCause(rfc2866.AcctTerminateCause_Value_ServiceUnavailable):
		return TerminateCauseServiceUnavailable
	case TerminateCause(rfc2866.AcctTerminateCause_Value_Callback):
		return TerminateCauseCallback
	case TerminateCause(rfc2866.AcctTerminateCause_Value_HostRequest):
		return TerminateCauseHostRequest
	}
	return ""
}

type SessionsLogData struct {
	CommonData
	AcctSessionId string    `gorm:"column:acct_session_id;unique;not null;index"`
	Login         string    `gorm:"column:login;not null"`
	StartTime     time.Time `json:"start_time"`
	StopTime      time.Time `json:"stop_time"`
	Duration      int64     `json:"duration"`
	NasName       string    `json:"nas_name"`
	Ip            string    `gorm:"column:ip;not null"`
	Mac           string    `json:"mac"`
	Circuit       string    `json:"circuit"`
	Remote        string    `json:"remote"`
	ServiceName   string    `json:"service_name"`
	LastAlive     time.Time `json:"last_alive"`
	OctetsIn      int64     `json:"octets_in"`
	OctetsOut     int64     `json:"octets_out"`

	CustomerId uuid.NullUUID `gorm:"column:customer_id;type:uuid;not null'"`
	Customer   CustomerData  `gorm:"foreignKey:CustomerId"`

	NasId   uuid.NullUUID `gorm:"column:nas_id;type:uuid"`
	Nas     NasData       `gorm:"foreignKey:NasId"`
	NasPort uint32        `json:"nas_port"`

	ServiceInternetId uuid.NullUUID       `gorm:"column:service_internet_id;type:uuid;not null"`
	ServiceInternet   ServiceInternetData `gorm:"foreignKey:ServiceInternetId"`

	TerminateCause string `gorm:"column:terminate_cause;default:'User-Request'"`
}

func (SessionsLogData) TableName() string {
	return "ncc_session_log"
}

func (SessionsLogData) FromSession(session SessionData) SessionsLogData {
	return SessionsLogData{
		CommonData:        session.CommonData,
		AcctSessionId:     session.AcctSessionId,
		Login:             session.Login,
		StartTime:         session.StartTime,
		StopTime:          time.Now(),
		Duration:          session.Duration,
		NasName:           session.NasName,
		Ip:                session.Ip,
		Mac:               session.Mac,
		Circuit:           session.Circuit,
		Remote:            session.Remote,
		ServiceName:       session.ServiceName,
		LastAlive:         session.LastAlive,
		OctetsIn:          session.OctetsIn,
		OctetsOut:         session.OctetsOut,
		CustomerId:        session.CustomerId,
		NasId:             session.NasId,
		ServiceInternetId: session.ServiceInternetId,
	}
}
