package models

import "github.com/google/uuid"

type HardwareModelData struct {
	CommonData
	Name                 string        `json:"name"`
	VendorId             uuid.NullUUID `gorm:"column:vendor_id;type:uuid;not null"`
	Vendor               VendorData    `gorm:"foreignKey:VendorId"`
	SnmpOidFdb           string        `json:"snmp_oid_fdb"`
	SnmpOidIfaceAdmState string        `json:"snmp_oid_iface_adm_state"`
	CidReverse           bool          `json:"cid_reverse"`
	CidPortOffset        int           `json:"cid_port_offset"`
	CidPortLen           int           `json:"cid_port_len"`
	CidVidOffset         int           `json:"cid_vid_offset"`
	CidVidLen            int           `json:"cid_vid_len"`
	Type                 int           `json:"type"`
}

func (HardwareModelData) TableName() string {
	return "ncc_hardware_model"
}
