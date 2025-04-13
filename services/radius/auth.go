package radius

import (
	"bytes"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"encoding/binary"
	"fmt"
	"github.com/gogf/gf/net/gipv4"
	"io"
	"net"
	"strconv"
	"strings"
)

type Limits struct {
	Login       string
	Password    string
	IdleTimeout uint8
	Ip          string
	SpeedIn     uint32
	SpeedOut    uint32
	BurstIn     uint32
	BurstOut    uint32
}

func (s *RadiusServer) passwordAuth(req *Packet, limits Limits) ([]AttrEncoder, error) {
	var reply []AttrEncoder

	if strings.ToLower(string(req.Attr(UserName))) == strings.ToLower(limits.Login) {

		if limits.Password == "" {
			return nil, fmt.Errorf("no password provided")
		}

		if req.HasAttr(UserPassword) {
			pass, err := DecryptPassword(req.Attr(UserPassword), req)
			if err != nil {
				return nil, fmt.Errorf("can't decode password: %w", err)
			}
			if pass != limits.Password {
				return nil, fmt.Errorf("invalid password")
			}
			s.log.Trace("PAP login user=", limits.Login)
		} else if req.HasAttr(CHAPPassword) {
			challenge := req.Attr(CHAPChallenge)
			hash := req.Attr(CHAPPassword)

			if challenge == nil && req.Auth != nil {
				challenge = req.Auth
			}

			if challenge != nil && !CHAPMatch(limits.Password, hash, challenge) {
				return nil, fmt.Errorf("invalid password")
			}
		} else {
			// Search for MSCHAP attrs
			attrs := make(map[AttributeType]AttrEncoder)
			for _, attr := range req.Attrs {
				if attr.Type() == VendorSpecific {
					hdr := VendorSpecificHeader(attr.Bytes())
					if hdr.VendorId == Microsoft {
						attrs[AttributeType(hdr.VendorType)] = attr
					}
				}
			}

			if len(attrs) > 0 && len(attrs) != 2 {
				return nil, fmt.Errorf("MSCHAP: Missing attrs? MS-CHAP-Challenge/MS-CHAP-Response")
			} else if len(attrs) == 2 {
				// Collect our data
				challenge := DecodeChallenge(attrs[MSCHAPChallenge].Bytes()).Value
				if _, isV1 := attrs[MSCHAPResponse]; isV1 {
					// MSCHAPv1
					res := DecodeResponse(attrs[MSCHAPResponse].Bytes())
					if res.Flags == 0 {
						// If it is zero, the NT-Response field MUST be ignored and
						// the LM-Response field used.
						return nil, fmt.Errorf("MSCHAPv1: LM-Response not supported.")
					}
					if bytes.Compare(res.LMResponse, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}) != 0 {
						s.log.Debug("Sending Access-Reject: MSCHAPv1: LM-Response set. user=", limits.Login)
						//return
					}

					// Check for correctness
					calc, e := Encryptv1(challenge, limits.Password)
					if e != nil {
						return nil, fmt.Errorf("MSCHAPv1: Server-side processing error")
					}
					mppe, e := Mppev1(limits.Password)
					if e != nil {
						s.log.Debug("MPPEv1: " + e.Error())
						return nil, fmt.Errorf("MPPEv1: Server-side processing error")
					}

					if bytes.Compare(res.NTResponse, calc) != 0 {
						s.log.Debug("MSCHAPv1 user=", limits.Login, "mismatch expect=", calc, ", received=", res.NTResponse)
						return nil, fmt.Errorf("invalid password")
					}

					reply = append(reply, VendorAttr{
						Type:     VendorSpecific,
						VendorId: Microsoft,
						Values: []VendorAttrString{
							/* 1 Encryption-Allowed, 2 Encryption-Required */
							{
								Type:  MSMPPEEncryptionPolicy,
								Value: []byte{0x0, 0x0, 0x0, 0x01},
							},
						},
					}.Encode())

					reply = append(reply, VendorAttr{
						Type:     VendorSpecific,
						VendorId: Microsoft,
						Values: []VendorAttrString{
							/* encryption types, allow RC4[40/128bit] */
							{
								Type:  MSMPPEEncryptionTypes,
								Value: []byte{0x0, 0x0, 0x0, 0x06},
							},
						},
					}.Encode())

					reply = append(reply, VendorAttr{
						Type:     VendorSpecific,
						VendorId: Microsoft,
						Values: []VendorAttrString{
							/* mppe - encryption negotation key */
							{
								Type:  MSCHAPMPPEKeys,
								Value: mppe,
							},
						},
					}.Encode())

				} else if _, isV2 := attrs[MSCHAP2Response]; isV2 {
					// MSCHAPv2
					res := DecodeResponse2(attrs[MSCHAP2Response].Bytes())
					// todo: strange bug with some clients
					if res.Flags != 0 {
						s.log.Debug("WARN: MSCHAPv2: Flags should be set to 0 user=", limits.Login)
						//return
					}
					enc, e := Encryptv2(challenge, res.PeerChallenge, limits.Login, limits.Password)
					if e != nil {
						s.log.Debug("MSCHAPv2: " + e.Error())
						return nil, fmt.Errorf("MSCHAPv2: Server-side processing error")
					}
					send, recv := Mmpev2(req.Secret(), limits.Password, req.Auth, res.Response)

					if bytes.Compare(res.Response, enc.ChallengeResponse) != 0 {
						s.log.Debug("MSCHAPv2 user=", limits.Login, "mismatch expect=", enc.ChallengeResponse, "received=", res.Response)
						return nil, fmt.Errorf("invalid password")
					}
					s.log.Trace("MSCHAPv2 login user=", limits.Login)

					reply = append(reply, VendorAttr{
						Type:     VendorSpecific,
						VendorId: Microsoft,
						Values: []VendorAttrString{
							/* success challenge */
							{
								Type:  MSCHAP2Success,
								Value: append([]byte{res.Ident}, []byte(enc.AuthenticatorResponse)...),
							},
						},
					}.Encode())

					reply = append(reply, VendorAttr{
						Type:     VendorSpecific,
						VendorId: Microsoft,
						Values: []VendorAttrString{
							/* Recv-Key */
							{
								Type:  MSMPPERecvKey,
								Value: recv,
							},
						},
					}.Encode())

					reply = append(reply, VendorAttr{
						Type:     VendorSpecific,
						VendorId: Microsoft,
						Values: []VendorAttrString{
							/* Send-Key */
							{
								Type:  MSMPPESendKey,
								Value: send,
							},
						},
					}.Encode())

					reply = append(reply, VendorAttr{
						Type:     VendorSpecific,
						VendorId: Microsoft,
						Values: []VendorAttrString{
							/* 1 Encryption-Allowed, 2 Encryption-Required */
							{
								Type:  MSMPPEEncryptionPolicy,
								Value: []byte{0x0, 0x0, 0x0, 0x01},
							},
						},
					}.Encode())

					reply = append(reply, VendorAttr{
						Type:     VendorSpecific,
						VendorId: Microsoft,
						Values: []VendorAttrString{
							/* encryption types, allow RC4[40/128bit] */
							{
								Type:  MSMPPEEncryptionTypes,
								Value: []byte{0x0, 0x0, 0x0, 0x06},
							},
						},
					}.Encode())

				} else {
					return nil, fmt.Errorf("MSCHAP: Response1/2 not found")
				}
			}
		}
	}

	return reply, nil
}

func (s *RadiusServer) Auth(w io.Writer, req *Packet) {
	if e := ValidateAuthRequest(req); e != "" {
		s.log.Warn("auth.begin e=%v", e)
		return
	}

	login := string(req.Attr(UserName))

	nasIP := gipv4.Long2ip(binary.BigEndian.Uint32(req.Attr(NASIPAddress)))
	nasAttrs, err := s.getNasAttrs(nasIP)

	s.log.Debug("Auth: %d", req.Identifier)

	access := s.getAccess(login)
	if access == nil {
		s.accessReject(login, w, req, "No such user")
		return
	}

	if access.State != models.CustomerStateActive {
		s.accessReject(login, w, req, "Account blocked")
		return
	}

	limits := access.Limits

	if net.ParseIP(login) != nil {
		if access.HasLease {
			s.accessAccept(login, w, req, s.addAttributes([]AttrEncoder{}, nasAttrs, limits))
			return
		}
		s.accessReject(login, w, req, "Lease not found")
		return
	}

	reply, err := s.passwordAuth(req, limits)
	if err != nil {
		s.accessReject(login, w, req, err.Error())
		return
	}

	s.accessAccept(login, w, req, s.addAttributes(reply, nasAttrs, limits))
}

func (s *RadiusServer) isStringAttr(val string) bool {
	if s.strRegexp.MatchString(val) {
		return false
	}
	return true
}

func (s *RadiusServer) addAttributes(reply []AttrEncoder, attrs []NasAttr, limits Limits) []AttrEncoder {
	for _, attr := range attrs {
		var val = attr.Val

		val = strings.ReplaceAll(val, "${SPEED_IN}", fmt.Sprintf("%d", limits.SpeedIn))
		val = strings.ReplaceAll(val, "${SPEED_OUT}", fmt.Sprintf("%d", limits.SpeedOut))
		val = strings.ReplaceAll(val, "${SPEED_IN_BURST}", fmt.Sprintf("%d", limits.BurstIn))
		val = strings.ReplaceAll(val, "${SPEED_OUT_BURST}", fmt.Sprintf("%d", limits.BurstOut))
		val = strings.ReplaceAll(val, "${IP}", fmt.Sprintf("%s", limits.Ip))

		if attr.Vendor > 0 {
			reply = append(reply, VendorAttr{
				Type:     VendorSpecific,
				VendorId: attr.Vendor,
				Values: []VendorAttrString{
					{
						Type:  AttributeType(attr.Code),
						Value: []byte(val),
					},
				},
			}.Encode())
		} else {

			if s.isStringAttr(val) {
				reply = append(reply, NewAttr(
					AttributeType(attr.Code),
					[]byte(val),
					uint8(len(val)),
				))
			} else {
				val, _ := strconv.ParseUint(val, 10, 32)
				a := make([]byte, 4)
				binary.BigEndian.PutUint32(a, uint32(val))
				reply = append(reply, NewAttr(
					AttributeType(attr.Code),
					a,
					0,
				))
			}

		}
	}

	if limits.Ip != "" {
		if ip := net.ParseIP(limits.Ip); ip != nil {
			reply = append(reply, NewAttr(
				FramedIPAddress,
				net.ParseIP(limits.Ip).To4(),
				0,
			))
		} else {
			s.log.Debug("Invalid IP:", limits.Ip)
		}
	}

	return reply
}

func (s *RadiusServer) accessAccept(login string, w io.Writer, req *Packet, reply []AttrEncoder) {
	s.log.Info("Access-Accept: %s", login)
	_, err := w.Write(req.Response(AccessAccept, reply))
	if err != nil {
		s.log.Error("Can't write Accept: %v", err)
	}
}

func (s *RadiusServer) accessReject(login string, w io.Writer, req *Packet, msg string) {
	s.log.Info("Access-Reject: %s (%s)", login, msg)
	_, err := w.Write(DefaultPacket(req, AccessReject, msg))
	if err != nil {
		s.log.Error("Can't write Reject: %v", err)
	}
}
