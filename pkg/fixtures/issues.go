package fixtures

import "code.evixo.ru/ncc/ncc-backend/pkg/models"

const (
	IssueTypeFailure     = "38385fcd-8dc1-42cc-aab6-53dd97a6ae40"
	IssueTypeNewCustomer = "8eec5a9d-a006-4681-aa12-92c0ccf339a5"

	IssueUrgencyHigh   = "d7ee23f3-0eb9-47db-9199-8d5f1056fcbb"
	IssueUrgencyMedium = "c1b3449e-a324-48b8-bf70-b4cff19a1f10"
	IssueUrgencyLow    = "8b5633c0-802e-4c12-be96-83665de8b034"
)

func FakeIssueTypes() []models.IssueTypeData {
	return []models.IssueTypeData{
		{
			CommonData: models.CommonData{
				Id: IssueTypeFailure,
			},
			Name:        "Неполадки с интернетом",
			DefaultType: true,
			Color:       "ff0000",
		},
		{
			CommonData: models.CommonData{
				Id: IssueTypeNewCustomer,
			},
			Name:        "Новое подключение",
			DefaultType: false,
			Color:       "00ff00",
		},
	}
}

func FakeIssueUrgencies() []models.IssueUrgencyData {
	return []models.IssueUrgencyData{
		{
			CommonData: models.CommonData{
				Id: IssueUrgencyHigh,
			},
			Name:     "Высокий",
			Priority: 30,
		},
		{
			CommonData: models.CommonData{
				Id: IssueUrgencyMedium,
			},
			Name:     "Средний",
			Priority: 20,
		},
		{
			CommonData: models.CommonData{
				Id: IssueUrgencyLow,
			},
			Name:     "Низкий",
			Priority: 10,
		},
	}
}
