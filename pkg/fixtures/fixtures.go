package fixtures

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"database/sql"
	"github.com/labstack/gommon/random"
	"time"
)

var (
	Internet1 = FakeServicesInternet()[0]
)

func FakeCustomerByLogin(login string) *models2.CustomerData {
	for _, c := range FakeCustomers() {
		if c.Login == login {
			return &c
		}
	}

	return nil
}

func FakeInformings() []models2.InformingData {
	return []models2.InformingData{
		{
			CommonData: models2.CommonData{
				Id:        "2370eb9e-e52b-4616-bedf-ab1d32dd5ad7",
				CreateTs:  time.Now(),
				CreatedBy: "df378bff-41cb-41cb-80e7-546d56df7d64",
			},
			Message:   "Hello {login}! Your deposit is {deposit}",
			Type:      "",
			Descr:     "Test informing about deposit",
			Name:      "Deposit informing",
			Start:     time.Now().Add(-24 * time.Hour),
			State:     models2.InformingStateEnabled,
			Repeating: models2.InformingRepeatingDaily,
			Conditions: []models2.InformingConditionData{
				{
					CommonData: models2.CommonData{
						Id:        "fb070650-1544-48c6-aa31-254fee5a55a8",
						CreateTs:  time.Now(),
						CreatedBy: "0eec1d3c-d06e-45d9-b77b-7e50923baad8",
					},
					Field: models2.FieldDeposit,
					Expr:  models2.ExprLt,
					Val:   "50",
				},
				{
					CommonData: models2.CommonData{
						Id:        "c3740800-5f2e-4c98-b11d-aebc59f02b5e",
						CreateTs:  time.Now(),
						CreatedBy: "0eec1d3c-d06e-45d9-b77b-7e50923baad8",
					},
					Field: models2.FieldSent,
					Val:   models2.ExprFalse,
				},
				{
					CommonData: models2.CommonData{
						Id:        "a0de6802-21c7-4f7d-aad7-1237f0c2058a",
						CreateTs:  time.Now(),
						CreatedBy: "0eec1d3c-d06e-45d9-b77b-7e50923baad8",
					},
					Field: models2.FieldInternetState,
					Val:   models2.ExprFalse,
				},
				{
					CommonData: models2.CommonData{
						Id:        "94c4f7ee-b69d-4261-bc9d-e450dce3388f",
						CreateTs:  time.Now(),
						CreatedBy: "0eec1d3c-d06e-45d9-b77b-7e50923baad8",
					},
					Field: models2.FieldBlockingState,
					Val:   models2.ExprFalse,
				},
			},
		},
		{
			CommonData: models2.CommonData{
				Id:        "78f5ebbd-424e-4cf5-9dc7-5ad41113e2c9",
				CreateTs:  time.Now(),
				CreatedBy: "df378bff-41cb-41cb-80e7-546d56df7d64",
			},
			Message:   "Hello {login}!",
			Type:      "",
			Descr:     "Test custom informing",
			Name:      "Custom informing",
			Start:     time.Now().Add(-24 * time.Hour),
			State:     models2.InformingStateEnabled,
			Mode:      models2.InformingModeTest,
			Repeating: models2.InformingRepeatingNever,
			Conditions: []models2.InformingConditionData{
				{
					CommonData: models2.CommonData{
						Id:        "66a60c1c-282c-4cf3-8693-0c03443aae12",
						CreateTs:  time.Now(),
						CreatedBy: "0eec1d3c-d06e-45d9-b77b-7e50923baad8",
					},
					Field: models2.FieldVerified,
					Expr:  models2.ExprEq,
					Val:   "true",
				},
				{
					CommonData: models2.CommonData{
						Id:        "a6412860-ae16-4f4c-a002-45fe66b76086",
						CreateTs:  time.Now(),
						CreatedBy: "0eec1d3c-d06e-45d9-b77b-7e50923baad8",
					},
					Field: models2.FieldGroup,
					Expr:  models2.ExprEq,
					Val:   "testCustomersGroup",
				},
				{
					CommonData: models2.CommonData{
						Id:        "c970c7af-d05b-4d1f-a6e2-3de5fabda956",
						CreateTs:  time.Now(),
						CreatedBy: "0eec1d3c-d06e-45d9-b77b-7e50923baad8",
					},
					Field: models2.FieldSent,
					Val:   models2.ExprFalse,
				},
			},
		},
	}
}

func FakeInformingsTestCustomers() []models2.InformingTestCustomerData {
	return []models2.InformingTestCustomerData{
		{
			CommonData: models2.CommonData{
				Id:        "d0e82ef1-e2a9-4b7d-a358-b79440fd2453",
				CreateTs:  time.Now(),
				CreatedBy: "c9c101c7-6da5-4332-9931-db0daadffd3c",
			},
			Customer: models2.CustomerData{
				CommonData: models2.CommonData{
					Id:        "dfa1bd91-6901-430b-af69-a5b3b1d83be6",
					CreateTs:  time.Now(),
					CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
				},
				Group: models2.CustomerGroupData{
					Name: "testCustomersGroup",
				},
				Uid:                  random.New().String(12, "0123456789"),
				Login:                "test1",
				Password:             "testPassw0rd",
				Phone:                "70000100002",
				Deposit:              1000.0,
				Credit:               0.0,
				BlockingState:        models2.CustomerStateActive,
				ServiceInternetState: models2.ServiceStateEnabled,
				ServiceInternet:      Internet1,
				ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
				Flags: []models2.CustomerFlagData{
					{
						Name: models2.FieldSent,
						Val:  models2.ExprFalse,
					},
				},
				VerifiedTs: models2.NewNullTime(time.Now()),
			},
		},
	}
}

func FakeCustomers() []models2.CustomerData {
	return []models2.CustomerData{
		{ // снять АП, 10, стейт не менять
			CommonData: models2.CommonData{
				Id:        "b39de918-5c2e-4aa9-b171-ce9dbdfc0a39",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Group: models2.CustomerGroupData{
				Name: "testGroup",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user1",
			Password:             "testPassw0rd",
			Phone:                "70000000001",
			Deposit:              100.0,
			Credit:               0.0,
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			VerifiedTs:           models2.NewNullTime(time.Now()),
		},
		{ // пропустить
			CommonData: models2.CommonData{
				Id:        "dfa1bd91-6901-430b-af69-a5b3b1d83be6",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user2",
			Password:             "testPassw0rd",
			Phone:                "70000000002",
			Deposit:              0.0,
			Credit:               0.0,
			BlockingState:        models2.CustomerStateBlocked,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprTrue,
				},
			},
			VerifiedTs: models2.NewNullTime(time.Now()),
		},
		{ // снять АП, 10, стейт не менять
			CommonData: models2.CommonData{
				Id:        "cbb00027-f6a5-4703-91e6-5ede6239b0f8",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Group: models2.CustomerGroupData{
				Name: "testGroup",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user3",
			Password:             "testPassw0rd",
			Phone:                "70000000003",
			Deposit:              5.0,
			Credit:               0.0,
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprTrue,
				},
			},
		},
		{ // АП не снимать, заблокировать
			CommonData: models2.CommonData{
				Id:        "3d90cc95-a40e-4bb9-8fd6-d45ed51e45a7",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user4",
			Password:             "testPassw0rd",
			Phone:                "70000000004",
			Deposit:              -5.0,
			Credit:               0.0,
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprTrue,
				},
			},
		},
		{ // снять АП, 12 (+IP), стейт не менять
			CommonData: models2.CommonData{
				Id:        "6294ef4f-f2d0-4784-8c1e-a9e27b8410fc",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user5",
			Password:             "testPassw0rd",
			Phone:                "70000000005",
			Deposit:              100,
			Credit:               0.0,
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprTrue,
				},
			},
		},
		{ // АП не снимать, заблокировать
			CommonData: models2.CommonData{
				Id:        "07c7b866-956f-4c13-8db3-451f83afd00a",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user6",
			Password:             "testPassw0rd",
			Phone:                "70000000006",
			Deposit:              0.0003,
			Credit:               0.0,
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprTrue,
				},
			},
		},
		{ // снять АП, 10, уменьшить дни, стейт не менять
			CommonData: models2.CommonData{
				Id:        "e12042ce-22a1-4e6b-95f9-662c7c844322",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user7",
			Password:             "testPassw0rd",
			Phone:                "70000000007",
			Deposit:              -5.0,
			Credit:               0.0,
			CreditDaysLeft:       2,
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprTrue,
				},
			},
		},
		{ // снять АП, 10, уменьшить дни, стейт не менять
			CommonData: models2.CommonData{
				Id:        "696bb461-7a3b-457e-9ed1-27ced8ec4122",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user8",
			Password:             "testPassw0rd",
			Phone:                "70000000008",
			Deposit:              -5.0,
			Credit:               0.0,
			CreditDaysLeft:       1,
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprTrue,
				},
			},
		},
		{ // снять АП, 10, стейт не менять
			CommonData: models2.CommonData{
				Id:        "56285bc8-3740-4789-9e55-c4b2e0d04e90",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:      random.New().String(12, "0123456789"),
			Login:    "user9",
			Password: "testPassw0rd",
			Phone:    "70000000009",
			Deposit:  -5.0,
			Credit:   0.0,
			CreditExpire: sql.NullTime{
				Time:  time.Date(2022, 9, 10, 4, 0, 0, 0, time.Now().Location()),
				Valid: true,
			},
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprTrue,
				},
			},
		},
		{ // АП не снимать, заблокировать
			CommonData: models2.CommonData{
				Id:        "6ed09554-ded8-4869-b7cb-5d6d768c06e5",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:      random.New().String(12, "0123456789"),
			Login:    "user10",
			Password: "testPassw0rd",
			Phone:    "70000000010",
			Deposit:  -5.0,
			Credit:   0.0,
			CreditExpire: sql.NullTime{
				Time:  time.Date(2022, 9, 8, 2, 0, 0, 0, time.Now().Location()),
				Valid: true,
			},
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprTrue,
				},
			},
		},
		{ // снять АП 10, стейт не менять
			CommonData: models2.CommonData{
				Id:        "9651f1bd-7d6e-4a78-8a0c-48502211915a",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user11",
			Password:             "testPassw0rd",
			Phone:                "70000000011",
			Deposit:              -5.0,
			Credit:               10.0,
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprTrue,
				},
			},
		},
		{ // АП не снимать, заблокировать
			CommonData: models2.CommonData{
				Id:        "5244bf91-c859-4b46-b997-74eb8efa48fb",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user12",
			Password:             "testPassw0rd",
			Phone:                "70000000012",
			Deposit:              -15.0,
			Credit:               10.0,
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprTrue,
				},
			},
		},
		{ // АП не снимать, заблокировать
			CommonData: models2.CommonData{
				Id:        "f144b7b3-ae1b-47d6-8d1e-8598ffcce76c",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user13",
			Password:             "testPassw0rd",
			Phone:                "70000000013",
			Deposit:              0.0,
			Credit:               0.0,
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprTrue,
				},
			},
		},
		{ // снять АП, не блокировать
			CommonData: models2.CommonData{
				Id:        "0d24d5ef-2707-499d-a26b-03649865de4d",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user14",
			Password:             "testPassw0rd",
			Phone:                "70000000014",
			Deposit:              0.01,
			Credit:               0.0,
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprTrue,
				},
			},
		},
		{ // АП не снимать, заблокировать
			CommonData: models2.CommonData{
				Id:        "f6096a56-82b8-408d-a5ed-03001dd73517",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user15",
			Password:             "testPassw0rd",
			Phone:                "70000000015",
			Deposit:              0.001,
			Credit:               0.0,
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprTrue,
				},
			},
		},
		{
			CommonData: models2.CommonData{
				Id:        "669443ff-07bb-43f2-98da-83cb55f1e928",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Group: models2.CustomerGroupData{
				Name: "testGroup",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user16",
			Password:             "testPassw0rd",
			Phone:                "70000000016",
			Deposit:              45,
			Credit:               0.0,
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprFalse,
				},
			},
			VerifiedTs: models2.NewNullTime(time.Now()),
		},
		{
			CommonData: models2.CommonData{
				Id:        "b3c3542d-c329-434d-91bd-e9fb83ba2878",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user17",
			Password:             "testPassw0rd",
			Phone:                "70000000017",
			Deposit:              34,
			Credit:               0.0,
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
		},
		{
			CommonData: models2.CommonData{
				Id:        "34033b19-ef2b-4bfc-8112-b31d4824b7a8",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user18",
			Password:             "testPassw0rd",
			Phone:                "70000000018",
			Deposit:              42,
			Credit:               0.0,
			BlockingState:        models2.CustomerStateBlocked,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
		},
		{
			CommonData: models2.CommonData{
				Id:        "0a43b321-5d92-430e-8b30-7b005d6a0820",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user19",
			Password:             "testPassw0rd",
			Phone:                "70000000019",
			Deposit:              42,
			Credit:               0.0,
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprTrue,
				},
			},
		},
	}
}

func FakeServiceInternetCustomData() []models2.ServiceInternetCustomData {
	return []models2.ServiceInternetCustomData{
		{
			CustomerId: models2.NewNullUUID("6294ef4f-f2d0-4784-8c1e-a9e27b8410fc"), // user5
			Ip:         "10.10.30.45",                                               // pseudo-routable
		},
		{
			CustomerId: models2.NewNullUUID("e12042ce-22a1-4e6b-95f9-662c7c844322"), // user7
			Ip:         "172.18.134.12",                                             // internal
		},
		{
			CustomerId: models2.NewNullUUID("696bb461-7a3b-457e-9ed1-27ced8ec4122"), // user8
			Ip:         "172.18.0.12",                                               // internal
		},
		{
			CustomerId: models2.NewNullUUID("56285bc8-3740-4789-9e55-c4b2e0d04e90"), // user9
			Ip:         "172.22.1.35",                                               // not in any pool
		},
		{
			CustomerId: models2.NewNullUUID("cbb00027-f6a5-4703-91e6-5ede6239b0f8"), // user3
			Ip:         "bad IP",                                                    // bad IP
		},
	}
}

func FakeIpPools() []models2.IpPoolData {
	return []models2.IpPoolData{
		{
			Name:      "Free pool",
			PoolStart: "172.18.0.10",
			PoolEnd:   "172.18.255.255",
			Mask:      "255.255.0.0",
			Dns1:      "8.8.8.8",
			Dns2:      "1.1.1.1",
			IsPaid:    false,
		},
		{
			Name:      "Paid pool",
			PoolStart: "10.10.30.10",
			PoolEnd:   "10.10.30.255",
			Mask:      "255.255.255.0",
			Dns1:      "8.8.8.8",
			Dns2:      "1.1.1.1",
			IsPaid:    true,
		},
	}
}

func FakeFeeProcessedMap() map[string]models2.CustomerData {
	return map[string]models2.CustomerData{
		"user16": {
			CommonData: models2.CommonData{
				Id:        "669443ff-07bb-43f2-98da-83cb55f1e928",
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                "user16",
			Phone:                "70000000016",
			Deposit:              45,
			Credit:               0.0,
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetState: models2.ServiceStateEnabled,
			ServiceInternet:      Internet1,
			ServiceInternetId:    models2.NewNullUUID(Internet1.CommonData.Id),
			Flags: []models2.CustomerFlagData{
				{
					Name: models2.FieldSent,
					Val:  models2.ExprFalse,
				},
			},
		},
	}
}
