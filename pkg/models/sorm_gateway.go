package models

type SormGatewayData struct {
	CommonData

	Code     string `gorm:"column:code"`
	Name     string `gorm:"column:name"`
	Descr    string `gorm:"column:descr"`
	Country  string `gorm:"column:country"`
	Region   string `gorm:"column:region"`
	District string `gorm:"column:district"`
	City     string `gorm:"column:city"`
	Street   string `gorm:"column:street"`
	Build    string `gorm:"column:build"`
	Type     string `gorm:"column:type"`
	IP       string `gorm:"column:ip"`
}

func (SormGatewayData) TableName() string {
	return "ncc_sorm_gateway"
}
