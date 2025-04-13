package dto

type FakeDataClearUsecaseRequest struct {
	ClearCustomers bool `json:"clear_customers"`
	ClearGeo       bool `json:"clear_geo"`
	ClearDevices   bool `json:"clear_devices"`
}

func (FakeDataClearUsecaseRequest) Validate() error {
	return nil
}

func (FakeDataClearUsecaseRequest) FromFakeDataClearEndpointRequest(req FakeDataClearEndpointRequest) FakeDataClearUsecaseRequest {
	return FakeDataClearUsecaseRequest{
		ClearCustomers: req.ClearCustomers,
		ClearDevices:   req.ClearDevices,
		ClearGeo:       req.ClearGeo,
	}
}

func NewFakeDataClearUsecaseRequest() FakeDataClearUsecaseRequest {
	return FakeDataClearUsecaseRequest{}
}
