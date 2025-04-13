package dto

type RadiusKillSessionsEndpointRequest struct {
	Sessions int  `json:"sessions"`
	Random   bool `json:"random"`
}

func (RadiusKillSessionsEndpointRequest) Validate() error {
	return nil
}

func NewRadiusKillSessionsEndpointRequest() RadiusKillSessionsEndpointRequest {
	return RadiusKillSessionsEndpointRequest{}
}
