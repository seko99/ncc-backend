package models

import "github.com/google/uuid"

type RadiusVendorData struct {
	CommonData
	Name string `json:"name"`
	Code int    `json:"code"`
}

func (RadiusVendorData) TableName() string {
	return "ncc_radius_vendor"
}

type RadiusAttributeData struct {
	CommonData
	Name     string           `json:"name"`
	Code     int              `json:"code"`
	VendorId uuid.NullUUID    `gorm:"column:vendor_id;type:uuid;not null"`
	Vendor   RadiusVendorData `gorm:"foreignKey:VendorId"`
}

func (RadiusAttributeData) TableName() string {
	return "ncc_radius_attribute"
}

type RadiusAttributeLink struct {
	VendorId    uuid.NullUUID       `gorm:"column:vendor_id;type:uuid;not null;index:idx_vendor_attr,unique"`
	Vendor      RadiusVendorData    `gorm:"foreignKey:VendorId"`
	AttributeId uuid.NullUUID       `gorm:"column:attribute_id;type:uuid;not null;index:idx_vendor_attr,unique"`
	Attribute   RadiusAttributeData `gorm:"foreignKey:AttributeId"`
}

func (RadiusAttributeLink) TableName() string {
	return "ncc_radius_vendor_radius_attribute_link"
}
