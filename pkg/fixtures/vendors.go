package fixtures

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
)

const (
	FakeVendorDlink  = "34ce57a1-2da2-44ae-a7e5-005216ebd467"
	FakeVendorTplink = "4c62d4f9-923d-4351-8c39-6d43d9047326"
	FakeVendorEltex  = "f049d74b-2bc1-4582-98d4-e605575ead5d"
)

func FakeHardwareModels() []models2.HardwareModelData {
	return []models2.HardwareModelData{
		{
			CommonData: models2.CommonData{
				Id: "9f3f109b-d3fa-42b5-b12a-869bc11d80a9",
			},
			Name:     "DES-1210-28/ME/B3",
			VendorId: models2.NewNullUUID(FakeVendorDlink),
		},
	}
}

func FakeVendors() []models2.VendorData {
	return []models2.VendorData{
		{
			CommonData: models2.CommonData{
				Id: FakeVendorDlink,
			},
			Name: "D-Link",
		},
		{
			CommonData: models2.CommonData{
				Id: FakeVendorTplink,
			},
			Name: "TP-Link",
		},
		{
			CommonData: models2.CommonData{
				Id: FakeVendorEltex,
			},
			Name: "Eltex",
		},
	}
}
