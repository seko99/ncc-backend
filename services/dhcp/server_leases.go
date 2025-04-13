package dhcp

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"fmt"
	"github.com/gogf/gf/net/gipv4"
	"sync"
	"time"
)

type ServerLease struct {
	Tid              uint32
	Login            string
	Ip               uint32
	Subnet           uint32
	Router           uint32
	DNS              []uint32
	Mac              string
	Remote           string
	Hostname         string
	CVID             uint16
	Port             uint16
	Start            uint32
	LeaseTime        uint32
	Expire           uint32
	CircuitId        string
	Status           uint16
	MarkedForRemoval bool
}

type ServerLeaseMap struct {
	sync.RWMutex
	m map[uint32]ServerLease
}

func (ths *ServerLease) FromLeaseData(data models.LeaseData) ServerLease {
	return ServerLease{
		Login:  data.Customer.Login,
		Ip:     gipv4.Ip2long(data.Ip),
		Subnet: gipv4.Ip2long(data.Subnet),
		Router: gipv4.Ip2long(data.Router),
		DNS: []uint32{
			gipv4.Ip2long(data.Dns1),
			gipv4.Ip2long(data.Dns2),
		},
		Mac:              data.Mac,
		Remote:           data.Remote,
		Hostname:         data.Hostname,
		CVID:             uint16(data.Cvid),
		Port:             uint16(data.Port),
		Start:            uint32(data.Start.Unix()),
		LeaseTime:        uint32(data.LeaseTime),
		Expire:           uint32(data.Expire.Unix()),
		CircuitId:        data.Circuit,
		Status:           uint16(data.Status),
		MarkedForRemoval: data.MarkedForRemoval,
	}
}

func (ths *ServerLease) ToLeaseData() models.LeaseData {
	return models.LeaseData{
		Ip:               gipv4.Long2ip(ths.Ip),
		Subnet:           gipv4.Long2ip(ths.Subnet),
		Router:           gipv4.Long2ip(ths.Router),
		Dns1:             gipv4.Long2ip(ths.DNS[0]),
		Dns2:             gipv4.Long2ip(ths.DNS[1]),
		Mac:              ths.Mac,
		Cvid:             int(ths.CVID),
		Port:             int(ths.Port),
		Remote:           ths.Remote,
		Circuit:          ths.CircuitId,
		LeaseTime:        int64(ths.LeaseTime),
		Start:            time.Unix(int64(ths.Start), 0),
		Expire:           time.Unix(int64(ths.Expire), 0),
		Hostname:         ths.Hostname,
		Status:           int(ths.Status),
		MarkedForRemoval: ths.MarkedForRemoval,
	}
}

func (ths *ServerLeaseMap) Set(m map[string]models.LeaseData) {
	ths.Lock()
	defer ths.Unlock()

	for _, v := range m {
		serverLease := ServerLease{}
		lease := serverLease.FromLeaseData(v)
		if l, ok := ths.m[lease.Ip]; ok {
			l.MarkedForRemoval = lease.MarkedForRemoval
			ths.m[lease.Ip] = l
		} else {
			ths.m[lease.Ip] = lease
		}
	}
}

func (ths *ServerLeaseMap) Get() map[uint32]ServerLease {
	ths.RLock()
	defer ths.RUnlock()

	leases := map[uint32]ServerLease{}

	for k, v := range ths.m {
		leases[k] = v
	}

	return leases
}

func (ths *ServerLeaseMap) GetByIp(ip uint32) (*ServerLease, error) {
	ths.RLock()
	defer ths.RUnlock()

	if lease, ok := ths.m[ip]; ok {
		return &lease, nil
	}

	return nil, fmt.Errorf("lease not found")
}

func (ths *ServerLeaseMap) GetByPacket(pkt *Packet) (*ServerLease, error) {
	ths.RLock()
	defer ths.RUnlock()

	var lease ServerLease

	for _, l := range ths.m {
		if l.Mac == byte2mac(pkt.Packet.ClientMAC) &&
			l.Remote == pkt.RemoteId &&
			l.Port == pkt.Port &&
			l.CVID == pkt.Cvid {
			return &l, nil
		}
	}

	return &lease, fmt.Errorf("lease not found: %s", byte2mac(pkt.Packet.ClientMAC))
}

func (ths *ServerLeaseMap) Create(lease ServerLease) error {
	ths.Lock()
	defer ths.Unlock()

	if _, ok := ths.m[lease.Ip]; ok {
		return fmt.Errorf("lease exists")
	}

	ths.m[lease.Ip] = lease

	return nil
}

func (ths *ServerLeaseMap) Delete(lease ServerLease) error {
	ths.Lock()
	defer ths.Unlock()

	if _, ok := ths.m[lease.Ip]; !ok {
		return fmt.Errorf("lease not found")
	}

	delete(ths.m, lease.Ip)

	return nil
}

func (ths *ServerLeaseMap) Update(lease ServerLease) error {
	ths.Lock()
	defer ths.Unlock()

	if _, ok := ths.m[lease.Ip]; !ok {
		return fmt.Errorf("lease not found: %s", gipv4.Long2ip(lease.Ip))
	}

	ths.m[lease.Ip] = lease

	return nil
}

func (ths *ServerLeaseMap) IsAllocated(ip uint32) bool {
	ths.RLock()
	defer ths.RUnlock()

	if _, ok := ths.m[ip]; ok {
		return true
	}

	return false
}

func (ths *ServerLeaseMap) AllocateInPool(pool *models.DhcpPoolData, bindingMap *BindingMap) (uint32, error) {
	start := gipv4.Ip2long(pool.RangeStart)
	end := gipv4.Ip2long(pool.RangeEnd)

	for ip := start; ip <= end; ip++ {
		if ths.IsAllocated(ip) {
			continue
		}

		if bindingMap.HasBinding(ip) {
			continue
		}

		return ip, nil
	}
	return 0, fmt.Errorf("pool overloaded")
}
