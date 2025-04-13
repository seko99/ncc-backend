package fixtures

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
)

const (
	FakeServiceInternetLight  = "91acf60c-e6e0-4fe4-977f-4d833bb8ab76"
	FakeServiceInternetMedium = "fcd8b646-b31f-491d-a61d-1cf1e3df1f00"
	FakeServiceInternetHigh   = "7439162f-3f69-4785-94a3-7a69bfd60aac"
)

func FakeServicesInternet() []models2.ServiceInternetData {
	return []models2.ServiceInternetData{
		{
			CommonData: models2.CommonData{
				Id: FakeServiceInternetLight,
			},
			CommonServiceData: models2.CommonServiceData{
				Name:    "Light",
				Fee:     300,
				FeeType: models2.FeeTypeDaily,
			},
			SpeedIn:  models2.Speed30M,
			SpeedOut: models2.Speed30M,
			IpFee:    60,
		},
		{
			CommonData: models2.CommonData{
				Id: FakeServiceInternetMedium,
			},
			CommonServiceData: models2.CommonServiceData{
				Name:    "Medium",
				Fee:     500,
				FeeType: models2.FeeTypeDaily,
			},
			SpeedIn:  models2.Speed50M,
			SpeedOut: models2.Speed50M,
			IpFee:    50,
		},
		{
			CommonData: models2.CommonData{
				Id: FakeServiceInternetHigh,
			},
			CommonServiceData: models2.CommonServiceData{
				Name:    "High",
				Fee:     700,
				FeeType: models2.FeeTypeDaily,
			},
			SpeedIn:  models2.Speed100M,
			SpeedOut: models2.Speed100M,
			IpFee:    50,
		},
	}
}
