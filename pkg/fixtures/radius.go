package fixtures

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

const (
	RadiusVendorCisco   = "4595899e-31ef-4139-8872-cb5f8d7d6b8c"
	RadiusVendorRFC2865 = "6722f7cd-5d26-4d26-a308-938184f3be9f"

	RadiusAttrIdleTimeout = "8a3c0fe9-48fb-45ba-ae05-0da93f551578"
	RadiusAttrFramedPool  = "cd5f1715-20cb-4b88-8103-9665f8862bfa"
)

func FakeRadiusAttrs() []models.RadiusAttributeData {
	return []models.RadiusAttributeData{
		{
			CommonData: models.CommonData{
				Id: "eef3ea06-3777-4f74-aff9-adb70518046b",
			},
			Name:     "SSG-Service-Info",
			Code:     251,
			VendorId: models.NewNullUUID(RadiusVendorCisco),
		},
		{
			CommonData: models.CommonData{
				Id: "176033b6-365f-4eab-bc06-a19a70e83a8e",
			},
			Name:     "Cisco-AVpair",
			Code:     1,
			VendorId: models.NewNullUUID(RadiusVendorCisco),
		},
		{
			CommonData: models.CommonData{
				Id: "891b679f-94b8-4207-b4f9-c3c0240f7471",
			},
			Name:     "Acct-Interim-Interval",
			Code:     85,
			VendorId: models.NewNullUUID(RadiusVendorRFC2865),
		},
		{
			CommonData: models.CommonData{
				Id: RadiusAttrIdleTimeout,
			},
			Name:     "Idle-Timeout",
			Code:     28,
			VendorId: models.NewNullUUID(RadiusVendorRFC2865),
		},
		{
			CommonData: models.CommonData{
				Id: RadiusAttrFramedPool,
			},
			Name:     "Framed-Pool",
			Code:     88,
			VendorId: models.NewNullUUID(RadiusVendorRFC2865),
		},
	}
}

func FakeRadiusVendors() []models.RadiusVendorData {
	return []models.RadiusVendorData{
		{
			CommonData: models.CommonData{
				Id: RadiusVendorCisco,
			},
			Name: "Cisco",
			Code: 9,
		},
		{
			CommonData: models.CommonData{
				Id: RadiusVendorRFC2865,
			},
			Name: "RFC2865",
			Code: 0,
		},
	}
}
