package models

type VendorData struct {
	CommonData
	Name string `json:"name"`
}

func (VendorData) TableName() string {
	return "ncc_vendor"
}
