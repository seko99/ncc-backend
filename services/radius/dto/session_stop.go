package dto

type SessionStopRequest struct {
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

type SessionStopResponse struct {
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
