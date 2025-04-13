package fixtures

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
)

const (
	FakePaymentTypeCash     = "5fca853f-4885-4071-8488-fe2a8ced3f33"
	FakePaymentTypeTerminal = "00a6019d-d675-4544-bd9e-b278df29ac71"
	FakePaymentTypeCard     = "c38af901-f7cd-433a-8b2d-3329f5477c4d"
)

func FakePaymentSystems(userId, customerId string) []models2.PaymentSystemData {
	return []models2.PaymentSystemData{
		{
			CommonData: models2.CommonData{
				Id: "ef30cd29-a1f8-410f-b001-0016ec7b1d16",
			},
			Enabled:        true,
			Name:           "Терминалы SuperPay",
			Token:          "8eab3c35-ed38-41e6-8c92-73cf54f31049",
			TestMode:       false,
			PaymentTypeId:  models2.NewNullUUID(FakePaymentTypeTerminal),
			UserId:         models2.NewNullUUID(userId),
			TestCustomerId: models2.NewNullUUID(customerId),
		},
		{
			CommonData: models2.CommonData{
				Id: "014c4b38-f36e-470d-9299-aa8bc23ee1f9",
			},
			Enabled:        true,
			Name:           "ЮKassa",
			Token:          "7f6e1eec-1670-4d54-b34e-891c96427617",
			TestMode:       false,
			PaymentTypeId:  models2.NewNullUUID(FakePaymentTypeCard),
			UserId:         models2.NewNullUUID(userId),
			TestCustomerId: models2.NewNullUUID(customerId),
		},
		{
			CommonData: models2.CommonData{
				Id: "a473682d-a26b-4270-ab0d-04d5d8806a39",
			},
			Enabled:        true,
			Name:           "Robokassa",
			Token:          "98e15d0d-fbba-426b-854f-6a52fe2eb6ae",
			TestMode:       false,
			PaymentTypeId:  models2.NewNullUUID(FakePaymentTypeCard),
			UserId:         models2.NewNullUUID(userId),
			TestCustomerId: models2.NewNullUUID(customerId),
		},
	}
}

func FakePaymentTypes() []models2.PaymentTypeData {
	return []models2.PaymentTypeData{
		{
			CommonData: models2.CommonData{
				Id: FakePaymentTypeCash,
			},
			AccountType:   models2.AccountTypeCash,
			Name:          "Наличные",
			ManualEnabled: true,
		},
		{
			CommonData: models2.CommonData{
				Id: FakePaymentTypeTerminal,
			},
			AccountType: models2.AccountTypeCashless,
			Name:        "Терминал",
		},
		{
			CommonData: models2.CommonData{
				Id: FakePaymentTypeCard,
			},
			AccountType: models2.AccountTypeCashless,
			Name:        "Card",
		},
	}
}
