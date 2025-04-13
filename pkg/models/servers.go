package models

import (
	"github.com/google/uuid"
	"time"
)

type ServerData struct {
	CommonData
	Name              string
	Descr             string
	GroupId           uuid.NullUUID   `gorm:"type:uuid;column:group_id"`
	Group             ServerGroupData `gorm:"foreignKey:GroupId"`
	Status            int
	StatusUpdated     time.Time
	MonitoringEnabled bool
}

func (ServerData) TableName() string {
	return "ncc_server"
}

type ServerGroupData struct {
	CommonData
	Name string
}

func (ServerGroupData) TableName() string {
	return "ncc_server_group"
}
