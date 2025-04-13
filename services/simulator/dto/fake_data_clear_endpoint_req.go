package dto

type FakeDataClearEndpointRequest struct {
	ClearCustomers bool `json:"clear_customers"`
	ClearDevices   bool `json:"clear_devices"`
	ClearGeo       bool `json:"clear_geo"`
}

func (FakeDataClearEndpointRequest) Validate() error {
	return nil
}

func NewFakeDataClearEndpointRequest() FakeDataClearEndpointRequest {
	return FakeDataClearEndpointRequest{}
}
