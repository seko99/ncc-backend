package models

type SormCustomerServiceData struct {
	CommonData

	Hash           string `gorm:"column:hash"`
	OrgUnit        string `gorm:"column:org_unit"`
	Login          string `gorm:"column:login"`
	ContractNumber string `gorm:"column:contract_number"`
	ServiceCode    string `gorm:"column:service_code"`
	Start          string `gorm:"column:start"`
	End            string `gorm:"column:end"`
	CustomData     string `gorm:"column:custom_data"`
}

func (SormCustomerServiceData) TableName() string {
	return "ncc_sorm_customer_services"
}
