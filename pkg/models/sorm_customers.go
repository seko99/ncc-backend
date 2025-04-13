package models

type SormCustomersData struct {
	CommonData

	Hash                  string `gorm:"column:hash"`
	OrgUnit               string `gorm:"column:org_unit"`
	Login                 string `gorm:"column:login"`
	IP                    string `gorm:"column:ip"`
	Email                 string `gorm:"column:email"`
	Phone                 string `gorm:"column:phone"`
	MAC                   string `gorm:"column:mac"`
	ContractDate          string `gorm:"column:contract_date"`
	ContractNumber        string `gorm:"column:contract_number"`
	Status                string `gorm:"column:status"`
	Start                 string `gorm:"column:start"`
	End                   string `gorm:"column:end"`
	Type                  string `gorm:"column:type"`
	NameType              string `gorm:"column:name_type"`
	FirstName             string `gorm:"column:first_name"`
	MiddleName            string `gorm:"column:middle_name"`
	LastName              string `gorm:"column:last_name"`
	Name                  string `gorm:"column:name"`
	BirthDate             string `gorm:"column:birth_date"`
	DocType               string `gorm:"column:doc_type"`
	DocSerial             string `gorm:"column:doc_serial"`
	DocNumber             string `gorm:"column:doc_number"`
	DocIssuer             string `gorm:"column:doc_issuer"`
	Doc                   string `gorm:"column:doc"`
	DocCode               string `gorm:"column:doc_code"`
	Bank                  string `gorm:"column:bank"`
	BankAccount           string `gorm:"column:bank_account"`
	EnterpriseName        string `gorm:"column:enterprise_name"`
	EnterpriseINN         string `gorm:"column:enterprise_inn"`
	ContactName           string `gorm:"column:contact_name"`
	ContactPhone          string `gorm:"column:contact_phone"`
	EnterpriseBank        string `gorm:"column:enterprise_bank"`
	EnterpriseBankAccount string `gorm:"column:enterprise_bank_account"`
	RegistrationType      string `gorm:"column:registration_type"`
	RegistrationZip       string `gorm:"column:registration_zip"`
	RegistrationCountry   string `gorm:"column:registration_country"`
	RegistrationRegion    string `gorm:"column:registration_region"`
	RegistrationDistrict  string `gorm:"column:registration_district"`
	RegistrationCity      string `gorm:"column:registration_city"`
	RegistrationStreet    string `gorm:"column:registration_street"`
	RegistrationBuild     string `gorm:"column:registration_build"`
	RegistrationCorp      string `gorm:"column:registration_corp"`
	RegistrationFlat      string `gorm:"column:registration_flat"`
	RegistrationAddress   string `gorm:"column:registration_address"`
	DeviceType            string `gorm:"column:device_type"`
	DeviceZip             string `gorm:"column:device_zip"`
	DeviceCountry         string `gorm:"column:device_country"`
	DeviceRegion          string `gorm:"column:device_region"`
	DeviceDistrict        string `gorm:"column:device_district"`
	DeviceCity            string `gorm:"column:device_city"`
	DeviceStreet          string `gorm:"column:device_street"`
	DeviceBuild           string `gorm:"column:device_build"`
	DeviceCorp            string `gorm:"column:device_corp"`
	DeviceFlat            string `gorm:"column:device_flat"`
	DeviceAddress         string `gorm:"column:device_address"`
	PostType              string `gorm:"column:post_type"`
	PostZip               string `gorm:"column:post_zip"`
	PostCountry           string `gorm:"column:post_country"`
	PostRegion            string `gorm:"column:post_region"`
	PostDistrict          string `gorm:"column:post_district"`
	PostCity              string `gorm:"column:post_city"`
	PostStreet            string `gorm:"column:post_street"`
	PostBuild             string `gorm:"column:post_build"`
	PostCorp              string `gorm:"column:post_corp"`
	PostFlat              string `gorm:"column:post_flat"`
	PostAddress           string `gorm:"column:post_address"`
	BillType              string `gorm:"column:bill_type"`
	BillZip               string `gorm:"column:bill_zip"`
	BillCountry           string `gorm:"column:bill_country"`
	BillRegion            string `gorm:"column:bill_region"`
	BillDistrict          string `gorm:"column:bill_district"`
	BillCity              string `gorm:"column:bill_city"`
	BillStreet            string `gorm:"column:bill_street"`
	BillBuild             string `gorm:"column:bill_build"`
	BillCorp              string `gorm:"column:bill_corp"`
	BillFlat              string `gorm:"column:bill_flat"`
	BillAddress           string `gorm:"column:bill_address"`
}

func (SormCustomersData) TableName() string {
	return "ncc_sorm_customers"
}
