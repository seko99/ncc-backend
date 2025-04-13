package models

type StreetData struct {
	CommonData
	Name string `json:"name"`
}

func (StreetData) TableName() string {
	return "ncc_street"
}
