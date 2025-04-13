package fixtures

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
)

func FakeCities() []models2.CityData {
	return []models2.CityData{
		{
			CommonData: models2.CommonData{
				Id: "8d5892a0-6d37-4681-891b-25b3c55cd79c",
			},
			Name: "Ростов-на-Дону",
		},
	}
}

func FakeStreets() []models2.StreetData {
	return []models2.StreetData{
		{
			CommonData: models2.CommonData{
				Id: "00000000-0000-0000-0000-000000000000",
			},
			Name: "technical street",
		},
	}
}
