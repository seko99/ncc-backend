package dto

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"fmt"
)

type SessionUpdateRequest struct {
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

type SessionUpdateResponse struct {
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

func (SessionUpdateResponse) FromStartResponse(req *SessionStartResponse) SessionUpdateResponse {
	return SessionUpdateResponse{
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

func (SessionUpdateResponse) FromSession(session *models.SessionData) SessionUpdateResponse {
	return SessionUpdateResponse{
		Timestamp:        session.UpdateTs.Unix(),
		UserName:         session.Login,
		AcctSessionId:    session.AcctSessionId,
		NasIpAddress:     session.Nas.Ip,
		NasPort:          fmt.Sprintf("%d", session.NasPort),
		NasId:            session.NasId.UUID.String(),
		FramedIp:         session.Ip,
		AcctInputOctets:  uint32(session.OctetsIn),
		AcctOutputOctets: uint32(session.OctetsOut),
		AcctSessionTime:  session.Duration,
	}
}
