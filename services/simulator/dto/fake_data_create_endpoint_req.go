package dto

type FakeDataCreateEndpointRequest struct {
	MaxCustomers        int    `json:"max_customers"`
	MinBuildLevel       int    `json:"min_build_level"`
	MaxMapNodes         int    `json:"max_map_nodes"`
	LeftUpper           LatLng `json:"left_upper"`
	RightBottom         LatLng `json:"right_bottom"`
	CreateCustomers     bool   `json:"create_customers"`
	CreateContracts     bool   `json:"create_contracts"`
	CreatePayments      bool   `json:"create_payments"`
	CreateFees          bool   `json:"create_fees"`
	CreateGeo           bool   `json:"create_geo"`
	DistributeCustomers bool   `json:"distribute_customers"`
	CreateDevices       bool   `json:"create_devices"`
	CreateBindings      bool   `json:"create_bindings"`
	CreateLeases        bool   `json:"create_leases"`
	CreateSessions      bool   `json:"create_sessions"`
}

func (FakeDataCreateEndpointRequest) Validate() error {
	return nil
}

func NewFakeDataCreateEndpointRequest() FakeDataCreateEndpointRequest {
	return FakeDataCreateEndpointRequest{}
}
