package fixtures

import "code.evixo.ru/ncc/ncc-backend/pkg/models"

func FakeDhcpPools() []models.DhcpPoolData {
	return []models.DhcpPoolData{
		{
			CommonData: models.CommonData{
				Id: "7d14e7ad-6e56-421c-a377-044d2d41aadd",
			},
			Type:              models.DhcpPoolTypeShared,
			Name:              "Users pool",
			RangeStart:        "172.45.0.10",
			RangeEnd:          "172.45.255.255",
			Mask:              "255.255.0.0",
			Gateway:           "172.45.0.1",
			Dns1:              "8.8.8.8",
			Dns2:              "1.1.1.1",
			LeaseTime:         360,
			UserChangeAllowed: false,
		},
		{
			CommonData: models.CommonData{
				Id: "df788473-eb3a-4e80-9951-2cdbc72b72c5",
			},
			Type:              models.DhcpPoolTypeStatic,
			Name:              "Static pool",
			RangeStart:        "144.93.25.10",
			RangeEnd:          "144.93.25.254",
			Mask:              "255.255.255.0",
			Gateway:           "144.93.25.1",
			Dns1:              "8.8.8.8",
			Dns2:              "1.1.1.1",
			LeaseTime:         360,
			UserChangeAllowed: false,
		},
	}
}
