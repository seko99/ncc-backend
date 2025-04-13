package fixtures

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
)

func FakeCustomerGroups() []models2.CustomerGroupData {
	return []models2.CustomerGroupData{
		{
			CommonData: models2.CommonData{
				Id: "c725ace0-a67f-47c1-b117-96268f6d41b9",
			},
			Name: "все",
		},
	}
}
