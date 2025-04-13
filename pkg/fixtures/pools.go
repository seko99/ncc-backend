package fixtures

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
)

func FakePools() []models2.DhcpPoolData {
	return []models2.DhcpPoolData{
		{
			CommonData: models2.CommonData{
				Id: "97c6c4e9-316a-410d-9422-cc445f6bca1a",
			},
			Name:              "User pool",
			Type:              models2.DhcpPoolTypeStatic,
			RangeStart:        "172.39.0.10",
			RangeEnd:          "172.39.255.255",
			Mask:              "255.255.0.0",
			Gateway:           "172.39.0.1",
			Dns1:              "8.8.8.8",
			Dns2:              "4.4.4.4",
			LeaseTime:         600,
			UserChangeAllowed: false,
		},
		{
			CommonData: models2.CommonData{
				Id: "601e0e3f-e575-41d2-922f-89840ae16f50",
			},
			Name:              "User shared pool",
			Type:              models2.DhcpPoolTypeShared,
			RangeStart:        "172.45.0.10",
			RangeEnd:          "172.45.255.255",
			Mask:              "255.255.0.0",
			Gateway:           "172.45.0.1",
			Dns1:              "8.8.8.8",
			Dns2:              "4.4.4.4",
			LeaseTime:         600,
			UserChangeAllowed: false,
		},
	}
}
