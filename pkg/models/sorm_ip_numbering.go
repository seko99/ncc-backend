package models

type SormIpNumberingData struct {
	CommonData

	Code    string `gorm:"column:code"`
	Name    string `gorm:"column:name"`
	Network string `gorm:"column:network"`
	Mask    string `gorm:"column:mask"`
}

func (SormIpNumberingData) TableName() string {
	return "ncc_sorm_ip_numbering"
}
