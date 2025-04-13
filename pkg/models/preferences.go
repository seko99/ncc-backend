package models

type PreferencesData struct {
	CommonData
	Name         string  `gorm:"column:name;unique"`
	StringValue  string  `gorm:"column:string_value"`
	IntValue     int     `gorm:"column:int_value"`
	DoubleValue  float64 `gorm:"column:double_value"`
	BooleanValue bool    `gorm:"column:boolean_value"`
	EntityLink   string  `gorm:"column:entity_link"`
}

func (PreferencesData) TableName() string {
	return "ncc_preferences"
}
