package models

type NasTypeData struct {
	CommonData
	Name          string             `json:"name"`
	NasAttributes []NasAttributeData `gorm:"many2many:ncc_nas_type_nas_attribute_link;joinForeignKey:nas_type_id;joinReferences:nas_attribute_id"`
}

func (NasTypeData) TableName() string {
	return "ncc_nas_type"
}
