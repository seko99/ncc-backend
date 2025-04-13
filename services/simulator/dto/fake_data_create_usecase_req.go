package dto

type FakeDataCreateUsecaseRequest struct {
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

func (FakeDataCreateUsecaseRequest) Validate() error {
	return nil
}

func (FakeDataCreateUsecaseRequest) FromFakeDataCreateEndpointRequest(req FakeDataCreateEndpointRequest) FakeDataCreateUsecaseRequest {
	return FakeDataCreateUsecaseRequest{
		MaxCustomers:        req.MaxCustomers,
		MinBuildLevel:       req.MinBuildLevel,
		MaxMapNodes:         req.MaxMapNodes,
		LeftUpper:           req.LeftUpper,
		RightBottom:         req.RightBottom,
		CreateCustomers:     req.CreateCustomers,
		CreateContracts:     req.CreateContracts,
		CreatePayments:      req.CreatePayments,
		CreateFees:          req.CreateFees,
		CreateGeo:           req.CreateGeo,
		DistributeCustomers: req.DistributeCustomers,
		CreateDevices:       req.CreateDevices,
		CreateBindings:      req.CreateBindings,
		CreateLeases:        req.CreateLeases,
		CreateSessions:      req.CreateSessions,
	}
}

func NewFakeDataCreateUsecaseRequest() FakeDataCreateUsecaseRequest {
	return FakeDataCreateUsecaseRequest{
		MinBuildLevel: 2,
		MaxMapNodes:   50,
		LeftUpper: LatLng{
			Lat: 39.6963,
			Lng: 47.2832,
		},
		RightBottom: LatLng{
			Lat: 39.7271,
			Lng: 47.2961,
		},
	}
}
