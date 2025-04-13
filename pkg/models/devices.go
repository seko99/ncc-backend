package models

import (
	"github.com/google/uuid"
	"time"
)

const (
	DeviceStatusOK    = 1
	DeviceStatusFault = 0

	DeviceIfaceStatusUp   = 1
	DeviceIfaceStatusDown = 2

	DeviceIfaceTypeNormal   = 10
	DeviceIfaceTypeFaulty   = 20
	DeviceIfaceTypeUplink   = 30
	DeviceIfaceTypeDownlink = 40
	DeviceIfaceTypeReserved = 50

	DeviceIfaceSpeed10   = 10_000_000
	DeviceIfaceSpeed100  = 100_000_000
	DeviceIfaceSpeed1000 = 1000_000_000
)

type DeviceStateData struct {
	CommonData
	StatusSnmp int `json:"status_snmp"`
	StatusIcmp int `json:"status_icmp"`
	//StatusSsh     int       `json:"status_ssh"`
	StatusUpdated time.Time `gorm:"column:status_updated"`
}

func (DeviceStateData) TableName() string {
	return "ncc_device_state"
}

type DeviceGroupData struct {
	CommonData
	Name string `gorm:"column:name"`
}

func (DeviceGroupData) TableName() string {
	return "ncc_device_group"
}

type DeviceData struct {
	CommonData
	ModelId           uuid.NullUUID     `gorm:"column:model_id"`
	Model             HardwareModelData `gorm:"foreignKey:ModelId"`
	GroupId           uuid.NullUUID     `gorm:"column:group_id"`
	Group             DeviceGroupData   `gorm:"foreignKey:GroupId"`
	Ip                string            `gorm:"column:ip"`
	Remote            string            `gorm:"column:remote"`
	SnmpRoCommunity   string            `gorm:"column:snmp_ro_community"`
	SnmpRwCommunity   string            `gorm:"column:snmp_rw_community"`
	Login             string            `gorm:"column:login"`
	Password          string            `gorm:"column:password"`
	StreetId          uuid.NullUUID     `gorm:"column:street_id;type:uuid"`
	Street            StreetData        `gorm:"foreignKey:StreetId"`
	CityId            uuid.NullUUID     `gorm:"column:city_id;type:uuid"`
	City              CityData          `gorm:"foreignKey:CityId"`
	Build             string            `gorm:"column:build"`
	Entrance          string            `gorm:"column:entrance"`
	Lat               float64           `gorm:"column:lat"`
	Lng               float64           `gorm:"column:lng"`
	Descr             string            `gorm:"column:descr"`
	Svid              int               `gorm:"column:svid"`
	SnmpTimeout       int               `gorm:"column:snmp_timeout"`
	SnmpVersion       string            `gorm:"column:snmp_version"`
	SnmpFdbOid        string            `gorm:"column:snmp_fdb_oid"`
	SnmpPort          int               `gorm:"column:snmp_port"`
	SnmpSleepTime     int               `gorm:"column:snmp_sleeptime"`
	MonitoringEnabled bool              `gorm:"column:monitoring_enabled"`
	PingEnabled       bool              `gorm:"column:ping_enabled"`
	PlaceOnDashboard  bool              `gorm:"column:place_on_dashboard"`
	SerialNumber      string            `gorm:"column:serial_number"`
	MapNodeId         uuid.NullUUID     `gorm:"column:map_node_id"`
	MapNode           MapNodeData       `gorm:"foreignKey:MapNodeId"`
	MonAliveEnabled   bool              `gorm:"column:mon_alive_enabled"`
	MonPortsEnabled   bool              `gorm:"column:mon_ports_enabled"`
	MonMetricsEnabled bool              `gorm:"column:mon_metrics_enabled"`
	MonFdbEnabled     bool              `gorm:"column:mon_fdb_enabled"`
	MonAliveType      int               `gorm:"column:mon_alive_type"`
	DeviceStateId     uuid.NullUUID     `gorm:"column:device_state_id"`
	DeviceState       DeviceStateData   `gorm:"foreignKey:DeviceStateId"`
	StatusIcmp        int               `gorm:"column:status_icmp"`
	StatusSnmp        int               `gorm:"column:status_snmp"`
	StatusSsh         int               `gorm:"column:status_ssh"`
	StatusUpdated     time.Time         `gorm:"column:status_updated"`
	PortCount         int               `gorm:"column:port_count"`
}

func (DeviceData) TableName() string {
	return "ncc_device"
}

type IfaceStateData struct {
	CommonData
	Speed       uint64    `gorm:"column:speed"`
	AdminStatus int       `gorm:"column:admin_status"`
	OperStatus  int       `gorm:"column:oper_status"`
	LastStatus  int       `gorm:"column:last_status"`
	LastChange  time.Time `gorm:"column:last_change"`
}

func (IfaceStateData) TableName() string {
	return "ncc_iface_state"
}

type IfaceData struct {
	CommonData

	DeviceId uuid.NullUUID `gorm:"column:device;type:uuid"`
	Device   DeviceData    `gorm:"foreignKey:DeviceId"`

	ServerId     uuid.NullUUID  `gorm:"column:server;type:uuid"`
	Hash         string         `gorm:"column:hash"`
	LinkDeviceId uuid.NullUUID  `gorm:"column:link_device;type:uuid"`
	LinkDevice   DeviceData     `gorm:"foreignKey:LinkDeviceId"`
	LinkIfaceId  uuid.NullUUID  `gorm:"column:link_iface;type:uuid"`
	Iface        string         `gorm:"column:iface"`
	Port         int            `gorm:"column:port"`
	Descr        string         `gorm:"column:descr"`
	OnDashboard  bool           `gorm:"column:on_dashboard"`
	IfaceStateId uuid.NullUUID  `gorm:"column:iface_state_id;type:uuid"`
	IfaceState   IfaceStateData `gorm:"foreignKey:IfaceStateId"`
	Type         int            `gorm:"column:type"`
}

func (IfaceData) TableName() string {
	return "ncc_iface"
}

type FDBData struct {
	CommonData
	DeviceId uuid.NullUUID `gorm:"column:device_id;type:uuid;not null"`
	Device   DeviceData    `gorm:"foreignKey:DeviceId"`

	CustomerId uuid.NullUUID `gorm:"column:customer_id;type:uuid"`
	Customer   CustomerData  `gorm:"foreignKey:CustomerId"`

	IfaceId uuid.NullUUID `gorm:"column:iface_id;type:uuid"`
	Iface   IfaceData     `gorm:"foreignKey:IfaceId"`

	IP       string    `gorm:"column:ip"`
	Vlan     int       `gorm:"column:vlan"`
	Port     int       `gorm:"column:port"`
	Mac      string    `gorm:"column:mac"`
	LastSeen time.Time `gorm:"column:last_seen"`
	Hash     string    `gorm:"column:hash"`
}

func (FDBData) TableName() string {
	return "ncc_fdb"
}

type PonONUState struct {
	CommonData

	Level     float64 `gorm:"column:level"`
	Distance  int     `gorm:"column:distance"`
	RTT       int     `gorm:"column:rtt"`
	IfaceName string  `gorm:"column:iface_name"`
	Status    int     `gorm:"column:status"`
}

func (PonONUState) TableName() string {
	return "ncc_pon_onu_state"
}

type PonONUData struct {
	CommonData

	DeviceId     uuid.NullUUID   `gorm:"column:device_id;type:uuid"`
	Device       DeviceData      `gorm:"foreignKey:DeviceId"`
	Port         string          `gorm:"column:port"`
	Mac          string          `gorm:"column:mac"`
	BindingId    uuid.NullUUID   `gorm:"column:binding_id;type:uuid"`
	Binding      DhcpBindingData `gorm:"foreignKey:BindingId"`
	OnuStateId   uuid.NullUUID   `gorm:"column:onu_state_id;type:uuid;not null"`
	OnuState     PonONUState     `gorm:"foreignKey:OnuStateId"`
	InvNumber    string          `gorm:"column:inv_number"`
	Vendor       string          `gorm:"column:vendor"`
	InstallDate  time.Time       `gorm:"column:install_date"`
	SerialNumber string          `gorm:"column:serial_number"`
}

func (PonONUData) TableName() string {
	return "ncc_pon_onu"
}

type TriggerData struct {
	CommonData

	Name                string        `gorm:"column:name"`
	Descr               string        `gorm:"column:descr"`
	ServerId            uuid.NullUUID `gorm:"column:server_id;type:uuid"`
	DeviceId            uuid.NullUUID `gorm:"column:device_id;type:uuid"`
	Device              DeviceData    `gorm:"foreignKey:DeviceId"`
	IfaceId             uuid.NullUUID `gorm:"column:iface_id;type:uuid"`
	Iface               IfaceData     `gorm:"foreignKey:IfaceId"`
	LastChange          time.Time     `gorm:"column:last_change"`
	PrevStatus          int           `gorm:"column:prev_status"`
	Status              int           `gorm:"column:status"`
	Severity            int           `gorm:"column:severity"`
	DependsOn           uuid.NullUUID `gorm:"column:depends_on;type:uuid"`
	NotificationEnabled bool          `gorm:"column:notification_enabled"`
}

func (TriggerData) TableName() string {
	return "ncc_trigger"
}

type PingerStatusData struct {
	CommonData

	LastUpdate time.Time     `gorm:"column:last_update"`
	IP         string        `gorm:"column:ip"`
	LastRTT    int64         `gorm:"column:last_rtt"`
	AvgRTT     float64       `gorm:"column:avg_rtt"`
	TotalRTT   int64         `gorm:"column:total_rtt"`
	Lost       int64         `gorm:"column:lost"`
	Sent       int64         `gorm:"column:sent"`
	LostRatio  float64       `gorm:"column:lost_ratio"`
	CustomerId uuid.NullUUID `gorm:"column:customer_id;type:uuid"`
	Customer   CustomerData  `gorm:"foreignKey:CustomerId"`
	DeviceId   uuid.NullUUID `gorm:"column:device_id;type:uuid"`
	Device     DeviceData    `gorm:"foreignKey:DeviceId"`
}

func (PingerStatusData) TableName() string {
	return "ncc_pinger_status"
}
