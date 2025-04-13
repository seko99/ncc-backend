package models

import (
	"github.com/google/uuid"
	"time"
)

const (
	IssueStatusOpen       = 10
	IssueStatusClosed     = 20
	IssueStatusInProgress = 30
	IssueStatusAssigned   = 40
)

type IssueTypeData struct {
	CommonData

	Name        string `json:"name"`
	DefaultType bool   `json:"default_type"`
	Color       string `json:"color"`
}

func (IssueTypeData) TableName() string {
	return "ncc_issue_type"
}

type IssueUrgencyData struct {
	CommonData

	Name     string `json:"name"`
	Priority int    `json:"priority"`
}

func (IssueUrgencyData) TableName() string {
	return "ncc_issue_urgency"
}

type IssueData struct {
	CommonData

	IId         int              `gorm:"column:iid"`
	Date        time.Time        `json:"date"`
	IssueTypeId uuid.NullUUID    `gorm:"column:issue_type_id;type:uuid"`
	IssueType   IssueTypeData    `gorm:"foreignKey:IssueTypeId"`
	CustomerId  uuid.NullUUID    `gorm:"column:customer_id;type:uuid"`
	Customer    CustomerData     `gorm:"foreignKey:CustomerId"`
	Status      int              `json:"status"`
	Address     string           `json:"address"`
	UrgencyId   uuid.NullUUID    `gorm:"column:urgency_id;type:uuid"`
	Urgency     IssueUrgencyData `gorm:"foreignKey:UrgencyId"`
	StreetId    uuid.NullUUID    `gorm:"column:street_id;type:uuid"`
	Street      StreetData       `gorm:"foreignKey:StreetId"`
	Build       string           `json:"build"`
	Flat        string           `json:"flat"`
	Name        string           `json:"name"`
	Phone       string           `json:"phone"`
	Comments    string           `json:"comments"`
	CityId      uuid.NullUUID    `gorm:"column:city_id;type:uuid"`
	City        CityData         `gorm:"foreignKey:CityId"`
}

func (IssueData) TableName() string {
	return "ncc_issue"
}

type IssueActionData struct {
	CommonData

	IssueId    uuid.NullUUID `gorm:"column:issue;type:uuid;not null"`
	Issue      IssueData     `gorm:"foreignKey:IssueId"`
	Status     int           `json:"status"`
	PrevStatus int           `json:"prev_status"`
	Comments   string        `json:"comments"`
}

func (IssueActionData) TableName() string {
	return "ncc_issue_action"
}
