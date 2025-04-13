package dhcp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

var MagicCookie = [4]byte{0x63, 0x82, 0x53, 0x63}

const (
	MsgRequest = 0x01
	MsgReply   = 0x02

	HwTypeEthernet = 0x01
	HwAddrLen      = 0x06

	BootpFlagsUnicast   = 0x0000
	BootpFlagsBroadcast = 0x8000

	MsgTypeDiscover = 1
	MsgTypeOffer    = 2
	MsgTypeRequest  = 3
	MsgTypeDecline  = 4
	MsgTypeAck      = 5
	MsgTypeNak      = 6
	MsgTypeRelease  = 7
	MsgTypeInform   = 8
)

type Packet struct {
	raw        []byte
	Packet     rawPacket
	Options    rawOptions
	Cvid       uint16
	Port       uint16
	RemoteId   string
	DomainName string
	circuitId  []byte
}

type rawPacket struct {
	MsgType        byte
	HwType         byte
	HwAddrLen      byte
	Hops           byte
	Tid            uint32
	SecElapsed     uint16
	Flags          uint16
	Ciaddr         uint32
	Yaddr          uint32
	NextServer     uint32
	RelayAgent     uint32
	ClientMAC      [6]byte
	_              [10]byte
	ServerHostname [64]byte
	BootfileName   [128]byte
	MagicCookie    [4]byte
}

type rawVendorOption struct {
	Opt uint16
	Val uint32
}

type rawOptions struct {
	Opt1 struct {
		// Subnet mask
		opt        byte
		len        byte
		SubnetMask uint32
	}
	Opt3 struct {
		// Router
		opt    byte
		len    byte
		Router uint32
	}
	Opt6 struct {
		// DNS
		opt byte
		len byte
		DNS []uint32
	}
	Opt12 struct {
		// Hostname
		opt      byte
		len      byte
		Hostname []byte
	}
	Opt15 struct {
		// Domain Name
		opt        byte
		len        byte
		DomainName []byte
	}
	Opt43 struct {
		// Vendor Specific Info
		opt            byte
		len            byte
		VendorSpecific []byte
	}
	Opt50 struct {
		// Requested IP
		opt         byte
		len         byte
		RequestedIP uint32
	}
	Opt51 struct {
		// Lease Time
		opt       byte
		len       byte
		LeaseTime uint32
	}
	Opt53 struct {
		// Message Type
		opt  byte
		len  byte
		Type byte
	}
	Opt54 struct {
		// DHCP Server ID
		opt          byte
		len          byte
		DHCPServerID uint32
	}
	Opt55 struct {
		// Parameter Request List
		opt                  byte
		len                  byte
		SubnetMask           bool
		Router               bool
		VendorSpecificInfo   bool
		NetbiosNameServer    bool
		NetbiosNodeType      bool
		NetbiosScope         bool
		DNS                  bool
		StaticRoute          bool
		ClasslessStaticRoute bool
		PrivateStaticRoute   bool
	}
	Opt57 struct {
		// Max DHCP Message Size
		opt                byte
		len                byte
		MaxDHCPMessageSize []byte
	}
	Opt58 struct {
		// Renew Time Value
		opt            byte
		len            byte
		RenewTimeValue int
	}
	Opt60 struct {
		// Vendor Class Identifier
		opt                   byte
		len                   byte
		VendorClassIdentifier string
	}
	Opt61 struct {
		// Client Identifier
		opt    byte
		len    byte
		HwType byte
		MAC    string
	}
	Opt81 struct {
		// Client FQDN
		opt        byte
		len        byte
		Flags      byte
		AA_RR      byte
		PTR_RR     byte
		ClientName string
	}
	Opt82 struct {
		// Agent Information Option
		opt            byte
		len            byte
		AgentCircuitID struct {
			opt byte
			len byte
			val []byte
		}
		AgentRemoteID struct {
			opt byte
			len byte
			val []byte
		}
		VendorSpecificInfo struct {
			opt byte
			len byte
			val []byte
		}
	}
	Opt116 struct {
		opt           byte
		len           byte
		AutoConfigure byte
	}
	Opt255 byte
}

func mac2byte(mac string) [6]byte {
	b, err := hex.DecodeString(strings.ReplaceAll(mac, ":", ""))
	if err != nil {
		return [6]byte{0, 0, 0, 0, 0, 0}
	}

	result := [6]byte{}
	for i := 0; i < 6; i++ {
		result[i] = b[i]
	}
	return result
}

func byte2mac(mac [6]byte) string {
	return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5])
}

func bytes2mac(mac []byte) (string, error) {
	if len(mac) > 0 {
		return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", mac[0], mac[1], mac[2], mac[3], mac[4], mac[5]), nil
	} else {
		return "", fmt.Errorf("mac len == 0")
	}
}

func (ths *Packet) ParsePacket() error {
	var pkt rawPacket

	data := ths.raw[0:binary.Size(pkt)]

	reader := bytes.NewReader(data)
	if err := binary.Read(reader, binary.BigEndian, &pkt); err != nil {
		return fmt.Errorf("can't read binary: %w", err)
	}

	ths.Packet = pkt

	return nil
}

func (ths *Packet) ParseOptions() error {
	var options rawOptions

	buf := ths.raw

	for p := binary.Size(rawPacket{}); p < binary.Size(buf); p++ {
		opt := buf[p]
		if opt == 255 {
			ths.Options = options
			return nil
		}
		p++
		if p >= binary.Size(buf) {
			return fmt.Errorf("buffer overflow")
		}
		optLen := buf[p]
		p++
		switch opt {
		case 1:
			options.Opt1.opt = opt
			options.Opt1.len = optLen
			reader := bytes.NewReader(buf[p : p+int(optLen)])
			if err := binary.Read(reader, binary.BigEndian, &options.Opt1.SubnetMask); err != nil {
				return fmt.Errorf("error reading Opt1")
			}
		case 3:
			options.Opt3.opt = opt
			options.Opt3.len = optLen
			reader := bytes.NewReader(buf[p : p+int(optLen)])
			if err := binary.Read(reader, binary.BigEndian, &options.Opt3.Router); err != nil {
				return fmt.Errorf("error reading Opt3")
			}
		case 6:
			options.Opt6.opt = opt
			options.Opt6.len = optLen
			reader := bytes.NewReader(buf[p : p+int(optLen)])
			if err := binary.Read(reader, binary.BigEndian, &options.Opt6.DNS); err != nil {
				return fmt.Errorf("error reading Opt6")
			}
		case 12:
			options.Opt12.opt = opt
			options.Opt12.len = optLen
			options.Opt12.Hostname = buf[p : p+int(optLen)]
		case 15:
			options.Opt15.opt = opt
			options.Opt15.len = optLen
			options.Opt15.DomainName = buf[p : p+int(optLen)]
		case 43:
			options.Opt43.opt = opt
			options.Opt43.len = optLen
			reader := bytes.NewReader(buf[p : p+int(optLen)])
			if err := binary.Read(reader, binary.BigEndian, &options.Opt43.VendorSpecific); err != nil {
				return fmt.Errorf("error reading Opt43")
			}
		case 50:
			options.Opt50.opt = opt
			options.Opt50.len = optLen
			reader := bytes.NewReader(buf[p : p+int(optLen)])
			if err := binary.Read(reader, binary.BigEndian, &options.Opt50.RequestedIP); err != nil {
				return fmt.Errorf("error reading Opt50")
			}
		case 51:
			options.Opt51.opt = opt
			options.Opt51.len = optLen
			reader := bytes.NewReader(buf[p : p+int(optLen)])
			if err := binary.Read(reader, binary.BigEndian, &options.Opt51.LeaseTime); err != nil {
				return fmt.Errorf("error reading Opt51")
			}
		case 53:
			options.Opt53.opt = opt
			options.Opt53.len = optLen
			options.Opt53.Type = buf[p]
		case 54:
			options.Opt54.opt = opt
			options.Opt54.len = optLen
			reader := bytes.NewReader(buf[p : p+int(optLen)])
			if err := binary.Read(reader, binary.BigEndian, &options.Opt54.DHCPServerID); err != nil {
				return fmt.Errorf("error reading Opt54")
			}
		case 55:
			options.Opt55.opt = opt
			options.Opt55.len = optLen
			for i := p; i < int(optLen); i++ {
				switch i {
				case 1:
					options.Opt55.SubnetMask = true
				case 3:
					options.Opt55.Router = true
				case 6:
					options.Opt55.DNS = true
				case 33:
					options.Opt55.StaticRoute = true
				case 43:
					options.Opt55.VendorSpecificInfo = true
				case 44:
					options.Opt55.NetbiosNameServer = true
				case 46:
					options.Opt55.NetbiosNodeType = true
				case 47:
					options.Opt55.NetbiosScope = true
				case 121:
					options.Opt55.ClasslessStaticRoute = true
				case 249:
					options.Opt55.PrivateStaticRoute = true
				}
			}
			break
		case 57:
			options.Opt57.opt = opt
			options.Opt57.len = optLen
			options.Opt57.MaxDHCPMessageSize = buf[p : p+int(optLen)]
		case 60:
			options.Opt60.opt = opt
			options.Opt60.len = optLen
			options.Opt60.VendorClassIdentifier = string(buf[p : p+int(optLen)])
		case 61:
			options.Opt61.opt = opt
			options.Opt61.len = optLen
			options.Opt61.HwType = buf[p]
			options.Opt61.MAC = fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", buf[p+1], buf[p+2], buf[p+3], buf[p+4], buf[p+5], buf[p+6])
		case 81:
			options.Opt81.opt = opt
			options.Opt81.len = optLen
			options.Opt81.Flags = buf[p]
			options.Opt81.AA_RR = buf[p+1]
			options.Opt81.PTR_RR = buf[p+2]
			options.Opt81.ClientName = string(buf[p+3 : p+int(optLen)])
		case 82:
			options.Opt82.opt = opt
			options.Opt82.len = optLen
			for i := p; i < p+int(optLen)-1; i++ {
				subopt := buf[i]
				i++
				suboptLen := buf[i]
				i++
				switch subopt {
				case 1:
					options.Opt82.AgentCircuitID.opt = subopt
					options.Opt82.AgentCircuitID.len = suboptLen
					options.Opt82.AgentCircuitID.val = buf[i : i+int(suboptLen)]
				case 2:
					options.Opt82.AgentRemoteID.opt = subopt
					options.Opt82.AgentRemoteID.len = suboptLen
					options.Opt82.AgentRemoteID.val = buf[i : i+int(suboptLen)]
				case 9:
					options.Opt82.VendorSpecificInfo.opt = subopt
					options.Opt82.VendorSpecificInfo.len = suboptLen
					options.Opt82.VendorSpecificInfo.val = buf[i : i+int(suboptLen)]
				}
				i = i + int(suboptLen) - 1
			}
			break
		case 116:
			options.Opt116.opt = opt
			options.Opt116.len = optLen
			options.Opt116.AutoConfigure = buf[p]
		}
		p = p + int(optLen) - 1
	}
	return fmt.Errorf("no Opt255 (end) found")
}

func (ths *Packet) parseRemoteId(remoteId []byte) (string, error) {
	var remote = ""

	if len(remoteId) < 2 {
		return "", fmt.Errorf("remote.id len <2")
	}
	if len(remoteId) == 6 {
		remote = fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", remoteId[0], remoteId[1], remoteId[2], remoteId[3], remoteId[4], remoteId[5])
	} else if len(remoteId) == 8 {
		r := remoteId[2:]
		remote = fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", r[0], r[1], r[2], r[3], r[4], r[5])
	} else {
		r := string(remoteId[2:])
		t := [13]string{}
		for i, c := range r {
			t[i] = string(c)
		}
		remote = fmt.Sprintf("%s%s:%s%s:%s%s:%s%s:%s%s:%s%s", t[0], t[1], t[2], t[3], t[4], t[5], t[6], t[7], t[8], t[9], t[10], t[11])
	}

	return remote, nil
}

func (ths *Packet) getPort() (uint16, error) {
	var port uint16 = 0

	if len(ths.circuitId) == 6 {
		portBytes := ths.circuitId[5:]
		if len(portBytes) > 0 {
			port = uint16(portBytes[0])
		} else {
			return 0, fmt.Errorf("port.bytes <= 0")
		}
	} else if len(ths.circuitId) == 5 {
		portBytes := ths.circuitId[4:]
		if len(portBytes) > 0 {
			port = uint16(portBytes[0])
		} else {
			return 0, fmt.Errorf("port.bytes <= 0")
		}
	}

	return port, nil
}

func (ths *Packet) getCvid() (uint16, error) {
	var cvid uint16 = 0

	if len(ths.circuitId) == 6 {
		cvidBytes := ths.circuitId[2:4]
		if len(cvidBytes) > 0 {
			cvid = binary.BigEndian.Uint16(cvidBytes)
		} else {
			return 0, fmt.Errorf("cvid.bytes <= 0")
		}
	} else if len(ths.circuitId) == 5 {
		cvidBytes := ths.circuitId[0:2]
		if len(cvidBytes) > 0 {
			cvid = binary.BigEndian.Uint16(cvidBytes)
		} else {
			return 0, fmt.Errorf("cvid.bytes <= 0")
		}
	}

	return cvid, nil
}

func (ths *Packet) Discover(mac string) ([]byte, error) {
	var buf bytes.Buffer
	pkt := rawPacket{}
	opt := rawOptions{}

	pkt.MsgType = MsgTypeRequest
	pkt.Hops = 0
	pkt.ClientMAC = mac2byte(mac)
	pkt.MagicCookie = MagicCookie

	opt.Opt53.opt = 53
	opt.Opt53.len = 1
	opt.Opt53.Type = MsgTypeDiscover

	opt.Opt255 = 0xFF

	if err := binary.Write(&buf, binary.BigEndian, pkt); err != nil {
		return nil, fmt.Errorf("error writing packet header")
	}

	if err := binary.Write(&buf, binary.BigEndian, opt.Opt53); err != nil {
		return nil, fmt.Errorf("error writing Opt53")
	}

	if err := binary.Write(&buf, binary.BigEndian, opt.Opt255); err != nil {
		return nil, fmt.Errorf("error writing Opt255 (end)")
	}

	return buf.Bytes(), nil
}

func (ths *Packet) Request(mac string, ip uint32) ([]byte, error) {
	var buf bytes.Buffer
	pkt := rawPacket{}
	opt := rawOptions{}

	pkt.MsgType = MsgTypeRequest
	pkt.Hops = 0
	pkt.ClientMAC = mac2byte(mac)
	pkt.Ciaddr = ip
	pkt.MagicCookie = MagicCookie

	opt.Opt53.opt = 53
	opt.Opt53.len = 1
	opt.Opt53.Type = MsgTypeRequest

	opt.Opt255 = 0xFF

	if err := binary.Write(&buf, binary.BigEndian, pkt); err != nil {
		return nil, fmt.Errorf("error writing packet header")
	}

	if err := binary.Write(&buf, binary.BigEndian, opt.Opt53); err != nil {
		return nil, fmt.Errorf("error writing Opt53")
	}

	if err := binary.Write(&buf, binary.BigEndian, opt.Opt255); err != nil {
		return nil, fmt.Errorf("error writing Opt255 (end)")
	}

	return buf.Bytes(), nil
}

func (ths *Packet) Offer(lease *ServerLease) ([]byte, error) {
	var buf bytes.Buffer
	pkt := rawPacket{}
	opt := rawOptions{}

	pkt.MsgType = MsgReply
	pkt.Hops = 0
	pkt.Ciaddr = lease.Ip
	pkt.MagicCookie = MagicCookie

	opt.Opt53.opt = 53
	opt.Opt53.len = 1
	opt.Opt53.Type = MsgTypeOffer

	opt.Opt255 = 0xFF

	if err := binary.Write(&buf, binary.BigEndian, pkt); err != nil {
		return nil, fmt.Errorf("error writing packet header")
	}

	if err := binary.Write(&buf, binary.BigEndian, opt.Opt53); err != nil {
		return nil, fmt.Errorf("error writing Opt53")
	}

	if err := binary.Write(&buf, binary.BigEndian, opt.Opt255); err != nil {
		return nil, fmt.Errorf("error writing Opt255 (end)")
	}

	return buf.Bytes(), nil
}

func (ths *Packet) Ack(lease *ServerLease) ([]byte, error) {
	var buf bytes.Buffer
	pkt := rawPacket{}
	opt := rawOptions{}

	pkt.MsgType = MsgReply
	pkt.Hops = 0
	pkt.Ciaddr = lease.Ip
	pkt.MagicCookie = MagicCookie

	opt.Opt53.opt = 53
	opt.Opt53.len = 1
	opt.Opt53.Type = MsgTypeAck

	opt.Opt255 = 0xFF

	if err := binary.Write(&buf, binary.BigEndian, pkt); err != nil {
		return nil, fmt.Errorf("error writing packet header")
	}

	if err := binary.Write(&buf, binary.BigEndian, opt.Opt53); err != nil {
		return nil, fmt.Errorf("error writing Opt53")
	}

	if err := binary.Write(&buf, binary.BigEndian, opt.Opt255); err != nil {
		return nil, fmt.Errorf("error writing Opt255 (end)")
	}

	return buf.Bytes(), nil
}

func (ths *Packet) Nak() ([]byte, error) {
	var buf bytes.Buffer
	pkt := rawPacket{}
	opt := rawOptions{}

	pkt.MsgType = MsgReply
	pkt.Hops = 0
	pkt.MagicCookie = MagicCookie

	opt.Opt53.opt = 53
	opt.Opt53.len = 1
	opt.Opt53.Type = MsgTypeNak

	opt.Opt255 = 0xFF

	if err := binary.Write(&buf, binary.BigEndian, pkt); err != nil {
		return nil, fmt.Errorf("error writing packet header")
	}

	if err := binary.Write(&buf, binary.BigEndian, opt.Opt53); err != nil {
		return nil, fmt.Errorf("error writing Opt53")
	}

	if err := binary.Write(&buf, binary.BigEndian, opt.Opt255); err != nil {
		return nil, fmt.Errorf("error writing Opt255 (end)")
	}

	return buf.Bytes(), nil
}

func (ths *Packet) FromLease(packetType byte, lease ServerLease) ([]byte, error) {
	pkt := ths.Packet
	options := ths.Options
	var opt rawOptions

	pkt.MsgType = MsgReply
	pkt.Hops = 0
	pkt.Yaddr = lease.Ip
	pkt.NextServer = lease.Router

	opt.Opt1.opt = 1
	opt.Opt1.len = 4
	opt.Opt1.SubnetMask = lease.Subnet

	opt.Opt3.opt = 3
	opt.Opt3.len = 4
	opt.Opt3.Router = lease.Router

	opt.Opt6.opt = 6
	opt.Opt6.len = byte(len(lease.DNS) * 4)
	opt.Opt6.DNS = lease.DNS

	opt.Opt15.opt = 15
	opt.Opt15.DomainName = []byte(ths.DomainName)
	opt.Opt15.len = byte(len(opt.Opt15.DomainName))

	opt.Opt51.opt = 51
	opt.Opt51.len = 4
	opt.Opt51.LeaseTime = lease.LeaseTime

	opt.Opt53.opt = 53
	opt.Opt53.len = 1
	opt.Opt53.Type = packetType

	opt.Opt54.opt = 54
	opt.Opt54.len = 4
	opt.Opt54.DHCPServerID = lease.Router

	opt.Opt255 = 0xFF

	var buf bytes.Buffer

	if err := binary.Write(&buf, binary.BigEndian, pkt); err != nil {
		return nil, fmt.Errorf("error writing packet header")
	}

	if err := binary.Write(&buf, binary.BigEndian, opt.Opt53); err != nil {
		return nil, fmt.Errorf("error writing Opt53")
	}
	if err := binary.Write(&buf, binary.BigEndian, opt.Opt54); err != nil {
		return nil, fmt.Errorf("error writing Opt54")
	}
	if err := binary.Write(&buf, binary.BigEndian, opt.Opt51); err != nil {
		return nil, fmt.Errorf("error writing Opt51")
	}
	if err := binary.Write(&buf, binary.BigEndian, opt.Opt1); err != nil {
		return nil, fmt.Errorf("error writing Opt1")
	}
	if err := binary.Write(&buf, binary.BigEndian, opt.Opt15.opt); err != nil {
		return nil, fmt.Errorf("error writing Opt15")
	}
	if err := binary.Write(&buf, binary.BigEndian, opt.Opt15.len); err != nil {
		return nil, fmt.Errorf("error writing Opt15.len")
	}
	if err := binary.Write(&buf, binary.BigEndian, opt.Opt15.DomainName); err != nil {
		return nil, fmt.Errorf("error writing Opt15.domainName")
	}
	if err := binary.Write(&buf, binary.BigEndian, opt.Opt3); err != nil {
		return nil, fmt.Errorf("error writing Opt3")
	}
	if err := binary.Write(&buf, binary.BigEndian, opt.Opt6.opt); err != nil {
		return nil, fmt.Errorf("error writing Opt6")
	}
	if err := binary.Write(&buf, binary.BigEndian, opt.Opt6.len); err != nil {
		return nil, fmt.Errorf("error writing Opt6.len")
	}
	if err := binary.Write(&buf, binary.BigEndian, opt.Opt6.DNS); err != nil {
		return nil, fmt.Errorf("error writing Opt6.dns")
	}
	if packetType == MsgTypeOffer || packetType == MsgTypeAck {
		if options.Opt60.VendorClassIdentifier == "MSFT 5.0" {
			opt.Opt43.opt = 43
			opt.Opt43.VendorSpecific = []byte{0x01, 0x04, 0x00, 0x00, 0x00, 0x02, 0x02, 0x04, 0x00, 0x00, 0x00, 0x01}
			opt.Opt43.len = byte(len(opt.Opt43.VendorSpecific))

			if err := binary.Write(&buf, binary.BigEndian, opt.Opt43.opt); err != nil {
				return nil, fmt.Errorf("error writing Opt43")
			}
			if err := binary.Write(&buf, binary.BigEndian, opt.Opt43.len); err != nil {
				return nil, fmt.Errorf("error writing Opt43.len")
			}
			if err := binary.Write(&buf, binary.BigEndian, opt.Opt43.VendorSpecific); err != nil {
				return nil, fmt.Errorf("error writing Opt43.vendor")
			}
		}
	}
	if options.Opt82.opt == 82 {
		if err := binary.Write(&buf, binary.BigEndian, options.Opt82.opt); err != nil {
			return nil, fmt.Errorf("error writing Opt82")
		}
		if err := binary.Write(&buf, binary.BigEndian, options.Opt82.len); err != nil {
			return nil, fmt.Errorf("error writing Opt82.len")
		}
		if options.Opt82.AgentCircuitID.opt == 1 {
			if err := binary.Write(&buf, binary.BigEndian, options.Opt82.AgentCircuitID.opt); err != nil {
				return nil, fmt.Errorf("error writing Opt82.agent")
			}
			if err := binary.Write(&buf, binary.BigEndian, options.Opt82.AgentCircuitID.len); err != nil {
				return nil, fmt.Errorf("error writing Opt82.agent.len")
			}
			if err := binary.Write(&buf, binary.BigEndian, options.Opt82.AgentCircuitID.val); err != nil {
				return nil, fmt.Errorf("error writing Opt82.agent.val")
			}
		}
		if options.Opt82.AgentRemoteID.opt == 2 {
			if err := binary.Write(&buf, binary.BigEndian, options.Opt82.AgentRemoteID.opt); err != nil {
				return nil, fmt.Errorf("error writing Opt82.remote")
			}
			if err := binary.Write(&buf, binary.BigEndian, options.Opt82.AgentRemoteID.len); err != nil {
				return nil, fmt.Errorf("error writing Opt82.remote.len")
			}
			if err := binary.Write(&buf, binary.BigEndian, options.Opt82.AgentRemoteID.val); err != nil {
				return nil, fmt.Errorf("error writing Opt82.remote.val")
			}
		}
		if options.Opt82.VendorSpecificInfo.opt == 9 {
			if err := binary.Write(&buf, binary.BigEndian, options.Opt82.VendorSpecificInfo.opt); err != nil {
				return nil, fmt.Errorf("error writing Opt82.vendor")
			}
			if err := binary.Write(&buf, binary.BigEndian, options.Opt82.VendorSpecificInfo.len); err != nil {
				return nil, fmt.Errorf("error writing Opt82.vendor.len")
			}
			if err := binary.Write(&buf, binary.BigEndian, options.Opt82.VendorSpecificInfo.val); err != nil {
				return nil, fmt.Errorf("error writing Opt82.vendor.val")
			}
		}
	}
	if err := binary.Write(&buf, binary.BigEndian, opt.Opt255); err != nil {
		return nil, fmt.Errorf("error writing Opt255 (end)")
	}

	return buf.Bytes(), nil
}

func NewPacket(data []byte) (*Packet, error) {
	p := &Packet{
		raw: data,
		//todo: option for domain name
		DomainName: "domain.local",
	}

	err := p.ParsePacket()
	if err != nil {
		return nil, fmt.Errorf("can't parse packet: %w", err)
	}

	err = p.ParseOptions()
	if err != nil {
		return nil, fmt.Errorf("can't parse options: %w", err)
	}

	p.circuitId = p.Options.Opt82.AgentCircuitID.val
	p.Cvid, err = p.getCvid()
	p.Port, err = p.getPort()
	p.RemoteId, err = p.parseRemoteId(p.Options.Opt82.AgentRemoteID.val)

	return p, nil
}
