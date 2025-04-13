package fixtures

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"github.com/bxcodec/faker/v4"
	"time"
)

func FakeLeases() []models.LeaseData {
	return []models.LeaseData{
		{
			CommonData: models.CommonData{
				Id:       "26a77957-7204-4083-8df3-1fd7f369b183",
				CreateTs: time.Now(),
			},
			Ip:         "10.50.0.11",
			CustomerId: models.NewNullUUID("b39de918-5c2e-4aa9-b171-ce9dbdfc0a39"),
			Mac:        faker.MacAddress(),
			Customer: models.CustomerData{
				CommonData: models.CommonData{
					Id: "b39de918-5c2e-4aa9-b171-ce9dbdfc0a39",
				},
				ServiceInternetId: models.NewNullUUID(Internet1.Id),
			},
		},
		{
			CommonData: models.CommonData{
				Id:       "85106e07-82d3-4150-b39f-8c897d2bd29f",
				CreateTs: time.Now(),
			},
			Ip:         "10.50.0.12",
			CustomerId: models.NewNullUUID("cbb00027-f6a5-4703-91e6-5ede6239b0f8"),
			Customer: models.CustomerData{
				CommonData: models.CommonData{
					Id: "b39de918-5c2e-4aa9-b171-ce9dbdfc0a39",
				},
				ServiceInternetId: models.NewNullUUID(Internet1.Id),
			},
		},
		{
			CommonData: models.CommonData{
				Id:       "ca65961e-2bbb-4473-9adc-f1e8ad3f5ec4",
				CreateTs: time.Now(),
			},
			Ip:         "10.50.0.13",
			CustomerId: models.NewNullUUID("3d90cc95-a40e-4bb9-8fd6-d45ed51e45a7"),
			Customer: models.CustomerData{
				CommonData: models.CommonData{
					Id: "b39de918-5c2e-4aa9-b171-ce9dbdfc0a39",
				},
				ServiceInternetId: models.NewNullUUID(Internet1.Id),
			},
		},
		{
			CommonData: models.CommonData{
				Id:       "d9b1ee21-cc87-4741-b5ce-889ff86edb05",
				CreateTs: time.Now(),
			},
			Ip: "10.50.0.14",
		},
	}
}
