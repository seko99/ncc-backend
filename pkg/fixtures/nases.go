package fixtures

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
)

const (
	NasTypeCisco = "4c9e894b-8100-43a7-a9fa-c9b9bbd7f2f2"
)

func FakeNasTypes() []models2.NasTypeData {
	return []models2.NasTypeData{
		{
			CommonData: models2.CommonData{
				Id: NasTypeCisco,
			},
			Name: "Cisco",
			NasAttributes: []models2.NasAttributeData{
				{
					AttributeId: RadiusAttrFramedPool,
					Val:         "og_pool",
				},
				{
					AttributeId: RadiusAttrIdleTimeout,
					Val:         "600",
				},
			},
		},
	}
}

func FakeNases() []models2.NasData {
	return []models2.NasData{
		{
			CommonData: models2.CommonData{
				Id: "05801662-cfff-4892-ab6f-d7e90afb70d2",
			},
			Name:            "bras00-central",
			Ip:              "127.0.0.1",
			Secret:          "someSecret902",
			SnmpCommunity:   "secretComm384",
			InterimInterval: 60,
			SessionTimeout:  180,
			DaeAddr:         "10.0.27.1",
			DaeSecret:       "superSecretPass932",
			DaeType:         models2.DaeTypeCoA,
			NasTypeId:       models2.NewNullUUID(NasTypeCisco),
		},
	}
}
