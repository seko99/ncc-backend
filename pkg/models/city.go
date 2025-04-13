package models

type CityData struct {
	CommonData
	Name string `json:"name"`
}

func (CityData) TableName() string {
	return "ncc_city"
}
