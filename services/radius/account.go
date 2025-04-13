package radius

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"encoding/binary"
	"io"
	"net"
	"time"
)

type SessionData struct {
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

func (s *RadiusServer) AcctOn(w io.Writer, req *Packet) {
	w.Write(DefaultPacket(req, AccountingResponse, "Updated accounting."))
}

func (s *RadiusServer) AcctOff(w io.Writer, req *Packet) {
	w.Write(DefaultPacket(req, AccountingResponse, "Updated accounting."))
}

func (s *RadiusServer) AcctBegin(w io.Writer, req *Packet) {
	if e := ValidateAcctRequest(req); e != "" {
		s.log.Warn("WARN: acct.begin err=", e)
		return
	}

	reply := []AttrEncoder{}
	user := ""
	framedIp := ""

	s.log.Trace("Start packet: %-v", req)

	if req.HasAttr(UserName) {
		user = string(req.Attr(UserName))
	} else {
		s.log.Debug("Empty User-Name, dropping start")
		w.Write(DefaultPacket(req, AccountingResponse, "Updated accounting."))
		return
	}

	if req.HasAttr(FramedIPAddress) {
		framedIp = DecodeIP(req.Attr(FramedIPAddress)).String()
	} else {
		if req.HasAttr(UserName) {
			if net.ParseIP(string(req.Attr(UserName))) != nil {
				framedIp = string(req.Attr(UserName))
			}
		}
	}

	sess := string(req.Attr(AcctSessionId))
	nasIp := DecodeIP(req.Attr(NASIPAddress)).String()
	//clientIp := string(req.Attr(CallingStationId))

	var nasIpAddress string
	if req.HasAttr(NASIPAddress) {
		nasIpAddress = DecodeIP(req.Attr(NASIPAddress)).String()
	}

	var serviceType uint32
	if req.HasAttr(ServiceType) {
		serviceType = binary.BigEndian.Uint32(req.Attr(ServiceType))
	}

	var callingStationId string
	var nasPort string

	if req.HasAttr(CallingStationId) {
		callingStationId = string(req.Attr(CallingStationId))
	}

	vendorAttrs, vendorAVPairs := getVendorAttrs(req)

	for _, av := range vendorAVPairs {
		if av.Attr == "connect-progress" {

		}
		if av.Attr == "client-mac-address" {
			callingStationId = av.Val
		}
	}

	if _, found := vendorAttrs[CiscoNASPort]; found {
		nasPort = vendorAttrs[CiscoNASPort].String()
	}

	sessionData := SessionData{
		Timestamp:        time.Now().UnixNano() / 1000,
		UserName:         user,
		NasIpAddress:     nasIpAddress,
		NasPort:          nasPort,
		AcctSessionId:    sess,
		FramedIp:         framedIp,
		CallingStationId: callingStationId,
		ServiceType:      serviceType,
	}

	s.log.Debug("acct.begin sessId=", sess, "USER=", user, "IP=", framedIp, "NAS=", nasIp)

	s.startSession(sessionData)

	w.Write(req.Response(AccountingResponse, reply))
}

func (s *RadiusServer) startSession(session SessionData) {
	err := s.events.PublishEvent(events.Event{
		Type: SessionStart,
		Payload: map[string]interface{}{
			"session": session,
		},
	})
	if err != nil {
		s.log.Error("Can't publish sessionStart: %v", err)
	}
}

func (s *RadiusServer) AcctUpdate(w io.Writer, req *Packet) {
	if e := ValidateAcctRequest(req); e != "" {
		s.log.Warn("acct.update e=", e)
		return
	}

	user := ""
	framedIp := ""

	if req.HasAttr(FramedIPAddress) {
		framedIp = DecodeIP(req.Attr(FramedIPAddress)).String()
	}

	if req.HasAttr(UserName) {
		user = string(req.Attr(UserName))
	} else {
		//user = framedIp
		s.log.Debug("Empty User-Name, dropping update")
		w.Write(DefaultPacket(req, AccountingResponse, "Updated accounting."))
		return
	}

	sess := string(req.Attr(AcctSessionId))
	nasIp := DecodeIP(req.Attr(NASIPAddress)).String()

	s.log.Debug("acct.update sessId=", sess, "USER=", user, "IP=", framedIp, "NAS=", nasIp, "sessTime=", req.Attr(AcctSessionTime))
	sessionId := string(req.Attr(AcctSessionId))
	sessTime := DecodeFour(req.Attr(AcctSessionTime))

	var octIn uint32
	var octOut uint32
	if req.HasAttr(AcctInputOctets) && req.HasAttr(AcctOutputOctets) {
		octIn = DecodeFour(req.Attr(AcctInputOctets))
		octOut = DecodeFour(req.Attr(AcctOutputOctets))
	}

	var nasIpAddress string
	if req.HasAttr(NASIPAddress) {
		nasIpAddress = DecodeIP(req.Attr(NASIPAddress)).String()
	}

	var serviceType uint32
	if req.HasAttr(ServiceType) {
		serviceType = binary.BigEndian.Uint32(req.Attr(ServiceType))
	}

	var callingStationId string
	var nasPort string

	if req.HasAttr(CallingStationId) {
		callingStationId = string(req.Attr(CallingStationId))
	}

	vendorAttrs, vendorAVPairs := getVendorAttrs(req)

	for _, av := range vendorAVPairs {
		if av.Attr == "connect-progress" {

		}
		if av.Attr == "client-mac-address" {
			callingStationId = av.Val
		}
	}

	if _, found := vendorAttrs[CiscoNASPort]; found {
		nasPort = vendorAttrs[CiscoNASPort].String()
	}

	sessionData := SessionData{
		Timestamp:        time.Now().UnixNano() / 1000,
		UserName:         user,
		NasIpAddress:     nasIpAddress,
		NasPort:          nasPort,
		AcctSessionId:    sessionId,
		FramedIp:         framedIp,
		CallingStationId: callingStationId,
		ServiceType:      serviceType,
		AcctSessionTime:  int64(sessTime),
		AcctInputOctets:  octIn,
		AcctOutputOctets: octOut,
	}

	s.updateSession(sessionData)

	w.Write(DefaultPacket(req, AccountingResponse, "Updated accounting."))
	//w.Write(DefaultPacket(req, AccountingResponse, ""))
}

func (s *RadiusServer) updateSession(session SessionData) {
	err := s.events.PublishEvent(events.Event{
		Type: InterimUpdate,
		Payload: map[string]interface{}{
			"session": session,
		},
	})
	if err != nil {
		s.log.Error("Can't publish interimUpdate: %v", err)
	}
}

func (s *RadiusServer) AcctStop(w io.Writer, req *Packet) {
	if e := ValidateAcctRequest(req); e != "" {
		s.log.Warn("acct.stop e=", e)
		return
	}

	user := ""
	framedIp := ""

	if req.HasAttr(FramedIPAddress) {
		framedIp = DecodeIP(req.Attr(FramedIPAddress)).String()
	}

	if req.HasAttr(UserName) {
		user = string(req.Attr(UserName))
	} else {
		user = framedIp
	}

	sess := string(req.Attr(AcctSessionId))
	nasIp := DecodeIP(req.Attr(NASIPAddress)).String()

	sessTime := DecodeFour(req.Attr(AcctSessionTime))

	var octIn uint32
	var octOut uint32
	if req.HasAttr(AcctInputOctets) && req.HasAttr(AcctOutputOctets) {
		octIn = DecodeFour(req.Attr(AcctInputOctets))
		octOut = DecodeFour(req.Attr(AcctOutputOctets))
	}

	var nasIpAddress string
	if req.HasAttr(NASIPAddress) {
		nasIpAddress = DecodeIP(req.Attr(NASIPAddress)).String()
	}

	var serviceType uint32
	if req.HasAttr(ServiceType) {
		serviceType = binary.BigEndian.Uint32(req.Attr(ServiceType))
	}

	var callingStationId string
	var nasPort string

	if req.HasAttr(CallingStationId) {
		callingStationId = string(req.Attr(CallingStationId))
	}

	vendorAttrs, vendorAVPairs := getVendorAttrs(req)

	for _, av := range vendorAVPairs {
		if av.Attr == "connect-progress" {

		}
		if av.Attr == "client-mac-address" {
			callingStationId = av.Val
		}
	}

	if _, found := vendorAttrs[CiscoNASPort]; found {
		nasPort = vendorAttrs[CiscoNASPort].String()
	}

	var acctTerminateCause uint32
	if req.HasAttr(AcctTerminateCause) {
		acctTerminateCause = DecodeFour(req.Attr(AcctTerminateCause))
	}

	sessionData := SessionData{
		Timestamp:          time.Now().UnixNano() / 1000,
		UserName:           user,
		NasIpAddress:       nasIpAddress,
		NasPort:            nasPort,
		AcctSessionId:      sess,
		FramedIp:           framedIp,
		CallingStationId:   callingStationId,
		ServiceType:        serviceType,
		AcctInputOctets:    octIn,
		AcctOutputOctets:   octOut,
		AcctTerminateCause: acctTerminateCause,
	}

	s.log.Debug("acct.stop sessId=", sess, "USER=", user, "IP=", framedIp, "NAS=", nasIp, "sessTime=", sessTime, "octetsIn=", octIn, "octetsOut=", octOut, "terminateCause=", acctTerminateCause)

	s.stopSession(sessionData)

	w.Write(DefaultPacket(req, AccountingResponse, "Finished accounting."))
}

func (s *RadiusServer) stopSession(session SessionData) {
	err := s.events.PublishEvent(events.Event{
		Type: SessionStop,
		Payload: map[string]interface{}{
			"session": session,
		},
	})
	if err != nil {
		s.log.Error("Can't publish sessionStop: %v", err)
	}
}
