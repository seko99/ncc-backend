package models

import "github.com/google/uuid"

const (
	DaeTypePOD  = 10
	DaeTypeCoA  = 20
	DaeTypeAuto = 30
)

type NasData struct {
	CommonData
	Name            string `json:"name"`
	Ip              string `json:"ip"`
	Secret          string `json:"secret"`
	SnmpCommunity   string `json:"snmp_community"`
	SessionTimeout  int32  `json:"session_timeout"`
	InterimInterval int32  `json:"interim_interval"`
	DaeAddr         string `json:"dae_addr"`
	DaeSecret       string `json:"dae_secret"`
	DaeType         int    `json:"dae_type"`

	NasTypeId uuid.NullUUID `gorm:"column:nas_type_id;type:uuid;not null"`
	NasType   NasTypeData   `gorm:"foreignKey:NasTypeId"`
}

func (NasData) TableName() string {
	return "ncc_nas"
}

type NasAttributeData struct {
	CommonData
	Attr        string              `json:"attr"`
	Val         string              `json:"val"`
	AttributeId string              `gorm:"column:attribute_id;type:uuid"`
	Attribute   RadiusAttributeData `gorm:"foreignKey:AttributeId"`
}

func (NasAttributeData) TableName() string {
	return "ncc_nas_attribute"
}

type NasTypeAttributeLink struct {
	NasTypeId      uuid.NullUUID    `gorm:"column:nas_type_id"`
	NasType        NasTypeData      `gorm:"foreignKey:NasTypeId"`
	NasAttributeId uuid.NullUUID    `gorm:"column:nas_attribute_id"`
	NasAttribute   NasAttributeData `gorm:"foreignKey:NasAttributeId"`
}

func (NasTypeAttributeLink) TableName() string {
	return "ncc_nas_type_nas_attribute_link"
}
