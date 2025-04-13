package radius

import (
	"encoding/binary"
	"strings"
)

type VendorAttrString struct {
	Type  AttributeType
	Value []byte
}

type VendorAttr struct {
	Type     AttributeType
	VendorId uint32
	Values   []VendorAttrString
}

type VendorHeader struct {
	VendorId   uint32
	VendorType uint8
}

type AVPair struct {
	Attr string
	Val  string
}

// Convert VendorAttr to generic Attr
func (t VendorAttr) Encode() AttrEncoder {
	val := make([]byte, 4)
	binary.BigEndian.PutUint32(val, t.VendorId)

	// Parse Values
	for _, value := range t.Values {
		raw := []byte(value.Value)

		b := make([]byte, 2+len(raw))
		b[0] = uint8(value.Type)   // vendor type
		b[1] = uint8(2 + len(raw)) // vendor length
		copy(b[2:], raw)
		//sum += 2+len(raw)
		val = append(val, b...)
	}

	return NewAttr(t.Type, val, 0)
}

func VendorSpecificHeader(b []byte) VendorHeader {
	return VendorHeader{
		VendorId:   binary.BigEndian.Uint32(b[0:4]),
		VendorType: b[4],
	}
}

func getVendorAttrs(req *Packet) (map[AttributeType]AttrEncoder, []AVPair) {
	vendorAttrs := make(map[AttributeType]AttrEncoder)
	var vendorAVPairs []AVPair
	for _, attr := range req.Attrs {
		if AttributeType(attr.Type()) == VendorSpecific {
			hdr := VendorSpecificHeader(attr.Bytes())
			if hdr.VendorId == Cisco {
				if attr.Type() == CiscoAVPair {
					av := strings.Split(attr.String(), "=")
					vendorAVPairs = append(vendorAVPairs, AVPair{Attr: av[0], Val: av[1]})
				} else {
					vendorAttrs[AttributeType(hdr.VendorType)] = attr
				}
			}
		}
	}
	return vendorAttrs, vendorAVPairs
}
