package models

type IpPoolData struct {
	CommonData

	Name              string `json:"name"`
	PoolStart         string `json:"pool_start"`
	PoolEnd           string `json:"pool_end"`
	OgPool            bool   `json:"og_pool"`
	UserChangeAllowed bool   `json:"user_change_allowed"`
	Mask              string `json:"mask"`
	Gateway           string `json:"gateway"`
	Dns1              string `json:"dns1"`
	Dns2              string `json:"dns2"`
	IsPaid            bool   `json:"is_paid"`
}

func (IpPoolData) TableName() string {
	return "ncc_ip_pool"
}
