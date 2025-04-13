package dto

type SessionStartRequest struct {
	Timestamp           int64
	UserName            string
	AcctSessionId       string
	NasIpAddress        string
	NasPort             string
	NasId               string
	NasPortType         string
	FramedIp            string
	FramedProtocol      string
	CallingStationId    string
	ServiceType         uint32
	EventTimestamp      string
	AcctAuthentic       string
	AcctInputOctets     uint32
	AcctOutputOctets    uint32
	AcctInputGigawords  uint8
	AcctOutputGigawords uint8
	AcctSessionTime     int64
	AcctTerminateCause  uint32
}

type SessionStartResponse struct {
	Timestamp           int64
	UserName            string
	AcctSessionId       string
	NasIpAddress        string
	NasPort             string
	NasId               string
	NasPortType         string
	FramedIp            string
	FramedProtocol      string
	CallingStationId    string
	ServiceType         uint32
	EventTimestamp      string
	AcctAuthentic       string
	AcctInputOctets     uint32
	AcctOutputOctets    uint32
	AcctInputGigawords  uint8
	AcctOutputGigawords uint8
	AcctSessionTime     int64
	AcctTerminateCause  uint32
}

func (SessionStartRequest) FromUpdateRequest(req SessionUpdateRequest) SessionStartRequest {
	return SessionStartRequest{
		Timestamp:           req.Timestamp,
		UserName:            req.UserName,
		AcctSessionId:       req.AcctSessionId,
		NasIpAddress:        req.NasIpAddress,
		NasPort:             req.NasPort,
		NasId:               req.NasId,
		NasPortType:         req.NasPortType,
		FramedIp:            req.FramedIp,
		FramedProtocol:      req.FramedProtocol,
		CallingStationId:    req.CallingStationId,
		ServiceType:         req.ServiceType,
		EventTimestamp:      req.EventTimestamp,
		AcctAuthentic:       req.AcctAuthentic,
		AcctInputOctets:     req.AcctInputOctets,
		AcctOutputOctets:    req.AcctOutputOctets,
		AcctInputGigawords:  req.AcctInputGigawords,
		AcctOutputGigawords: req.AcctOutputGigawords,
		AcctSessionTime:     req.AcctSessionTime,
		AcctTerminateCause:  req.AcctTerminateCause,
	}
}
