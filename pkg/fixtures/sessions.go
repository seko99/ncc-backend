package fixtures

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"github.com/bxcodec/faker/v4"
	"time"
)

func FakeSessions() []models.SessionData {
	return []models.SessionData{
		{
			CommonData: models.CommonData{
				Id:       "aa856c8f-c278-45a3-b6be-b3161fc8b5b3",
				CreateTs: time.Now(),
			},
			AcctSessionId: "53918348-f675-4412-926b-01a28e3a2639",
			CustomerId:    models.NewNullUUID("b39de918-5c2e-4aa9-b171-ce9dbdfc0a39"),
			Customer: models.CustomerData{
				CommonData: models.CommonData{
					Id: "b39de918-5c2e-4aa9-b171-ce9dbdfc0a39",
				},
				ServiceInternetId: models.NewNullUUID(Internet1.Id),
			},
			Login:             "user1",
			Ip:                "10.50.0.11",
			Mac:               faker.MacAddress(),
			StartTime:         time.Now(),
			OctetsIn:          1024_000,
			OctetsOut:         512_000,
			Duration:          135,
			ServiceInternetId: models.NewNullUUID(Internet1.Id),
			LastAlive:         time.Now().Add(-time.Second * 10),
			Nas:               FakeNases()[0],
		},
		{
			CommonData: models.CommonData{
				Id:       "2be2c35f-dc87-482f-9a6f-02ea26bdb7f4",
				CreateTs: time.Now(),
			},
			AcctSessionId: "3eead7bf-4851-402c-9236-cd7771ae6a4a",
			CustomerId:    models.NewNullUUID("cbb00027-f6a5-4703-91e6-5ede6239b0f8"),
			Customer: models.CustomerData{
				CommonData: models.CommonData{
					Id: "cbb00027-f6a5-4703-91e6-5ede6239b0f8",
				},
				ServiceInternetId: models.NewNullUUID(Internet1.Id),
			},
			Login:             "user3",
			Ip:                "10.50.0.12",
			StartTime:         time.Now(),
			OctetsIn:          1024_000,
			OctetsOut:         512_000,
			Duration:          135,
			ServiceInternetId: models.NewNullUUID(Internet1.Id),
			LastAlive:         time.Now().Add(-time.Second * 200),
		},
	}
}
