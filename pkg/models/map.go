package models

import "github.com/google/uuid"

const (
	MapNodeStatusActive   = 10
	MapNodeStatusInactive = 20
)

type MapNodeData struct {
	CommonData
	Lat      float64       `json:"lat"`
	Lng      float64       `json:"lng"`
	CityId   uuid.NullUUID `gorm:"column:city_id;type:uuid;not null"`
	City     CityData      `gorm:"foreignKey:CityId"`
	StreetId uuid.NullUUID `gorm:"column:street_id;type:uuid;not null"`
	Street   StreetData    `gorm:"foreignKey:StreetId"`
	Build    string        `json:"build"`
	Entrance string        `json:"entrance"`
	Status   int           `json:"status"`
}

func (MapNodeData) TableName() string {
	return "ncc_map_node"
}
