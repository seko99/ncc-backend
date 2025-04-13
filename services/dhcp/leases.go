package dhcp

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"encoding/hex"
	"fmt"
	"github.com/gogf/gf/net/gipv4"
	"sync"
	"time"
)

const (
	LeaseStatusAllocated = 0
	LeaseStatusAccepted  = 1
)

type LeaseMap struct {
	sync.RWMutex
	m map[string]models.LeaseData
}

func (ths *LeaseMap) Set(m map[string]models.LeaseData) {
	ths.Lock()
	defer ths.Unlock()

	ths.m = map[string]models.LeaseData{}
	for k, v := range m {
		ths.m[k] = v
	}
}

func (ths *Server) acceptLease(lease ServerLease) (*ServerLease, error) {
	lease.Status = LeaseStatusAccepted

	err := ths.serverLeases.Update(lease)
	if err != nil {
		return nil, fmt.Errorf("can't update server lease: %w", err)
	}

	return &lease, nil
}

func (ths *Server) removeLease(lease ServerLease) error {
	err := ths.serverLeases.Delete(lease)
	if err != nil {
		return fmt.Errorf("can't remove lease: %w", err)
	}

	return nil
}

func (ths *Server) newLease(pkt *Packet) (*ServerLease, error) {
	lease := ServerLease{}

	var pool *models.DhcpPoolData

	ths.allocWg.Wait()
	ths.allocWg.Add(1)
	defer ths.allocWg.Done()

	binding, err := ths.bindingMap.GetByPacket(pkt)
	if err != nil {
		pool, err = ths.poolMap.GetByType(models.DhcpPoolTypeShared)
		if err != nil {
			return nil, fmt.Errorf("can't find shared pool: %w", err)
		}

		ip, err := ths.serverLeases.AllocateInPool(pool, &ths.bindingMap)
		if err != nil {
			return nil, fmt.Errorf("can't allocate IP: %w", err)
		}

		lease.Ip = ip
	} else {
		lease.Ip = gipv4.Ip2long(binding.Ip)
		lease.Login = binding.Customer.Login

		pool, err = ths.poolMap.GetById(binding.PoolId.UUID.String())
		if err != nil {
			return nil, fmt.Errorf("can't find pool for binding: %w", err)
		}
	}

	lease.Subnet = gipv4.Ip2long(pool.Mask)
	lease.Router = gipv4.Ip2long(pool.Gateway)
	lease.DNS = append(lease.DNS, gipv4.Ip2long(pool.Dns1))
	lease.DNS = append(lease.DNS, gipv4.Ip2long(pool.Dns2))
	lease.LeaseTime = uint32(pool.LeaseTime)

	lease.Mac = byte2mac(pkt.Packet.ClientMAC)
	lease.Remote = pkt.RemoteId
	lease.CVID = pkt.Cvid
	lease.Port = pkt.Port
	lease.Hostname = string(pkt.Options.Opt12.Hostname)
	lease.CircuitId = hex.EncodeToString(pkt.circuitId)
	lease.Start = uint32(time.Now().Unix())
	lease.Expire = lease.Start + lease.LeaseTime
	lease.Status = LeaseStatusAllocated

	err = ths.serverLeases.Create(lease)
	if err != nil {
		return nil, fmt.Errorf("can't add server lease: %w", err)
	}

	return &lease, nil
}

func (ths *Server) renewLease(lease ServerLease) (*ServerLease, error) {
	pool, err := ths.poolMap.GetByIp(gipv4.Long2ip(lease.Ip))
	if err != nil {
		return nil, fmt.Errorf("can't find pool: %w", err)
	}

	lease.Expire = uint32(time.Now().Unix()) + uint32(pool.LeaseTime)

	err = ths.serverLeases.Update(lease)
	if err != nil {
		return nil, fmt.Errorf("can't update server lease: %w", err)
	}

	return &lease, nil
}

func (ths *Server) IsAllocated(ip uint32) bool {
	return ths.serverLeases.IsAllocated(ip)
}

func (ths *Server) GetServerLeases() map[uint32]ServerLease {
	return ths.serverLeases.Get()
}

func (ths *Server) GetServerLeaseByIP(ip uint32) (*ServerLease, error) {
	return ths.serverLeases.GetByIp(ip)
}
