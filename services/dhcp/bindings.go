package dhcp

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"fmt"
	"github.com/gogf/gf/net/gipv4"
	"sync"
)

type BindingMap struct {
	sync.RWMutex
	m map[string]models.DhcpBindingData
}

func (ths *BindingMap) Set(m map[string]models.DhcpBindingData) {
	ths.Lock()
	defer ths.Unlock()

	ths.m = map[string]models.DhcpBindingData{}
	for k, v := range m {
		ths.m[k] = v
	}
}

func (ths *BindingMap) Get() map[string]models.DhcpBindingData {
	ths.RLock()
	defer ths.RUnlock()

	data := map[string]models.DhcpBindingData{}

	for k, v := range ths.m {
		data[k] = v
	}

	return data
}

func (ths *BindingMap) GetByPacket(pkt *Packet) (*models.DhcpBindingData, error) {
	ths.RLock()
	defer ths.RUnlock()

	for _, b := range ths.m {
		if len(b.Mac) > 0 && b.Mac != byte2mac(pkt.Packet.ClientMAC) {
			continue
		}

		if len(b.Remote) > 0 && b.Remote != pkt.RemoteId {
			continue
		}

		if b.Cvid > 0 && b.Cvid != int(pkt.Cvid) {
			continue
		}

		if b.Port > 0 && b.Port != int(pkt.Port) {
			continue
		}

		return &b, nil
	}

	return nil, fmt.Errorf("binding not found")
}

func (ths *BindingMap) HasBinding(ip uint32) bool {
	ths.RLock()
	defer ths.RUnlock()

	if _, ok := ths.m[gipv4.Long2ip(ip)]; ok {
		return true
	}

	return false
}
