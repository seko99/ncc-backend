package dhcp

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"fmt"
	"github.com/gogf/gf/net/gipv4"
	"sync"
)

type PoolMap struct {
	sync.RWMutex
	m map[string]models.DhcpPoolData
}

func (ths *PoolMap) Set(m map[string]models.DhcpPoolData) {
	ths.Lock()
	defer ths.Unlock()

	ths.m = map[string]models.DhcpPoolData{}
	for k, v := range m {
		ths.m[k] = v
	}
}

func (ths *PoolMap) Get() map[string]models.DhcpPoolData {
	ths.RLock()
	defer ths.RUnlock()

	data := map[string]models.DhcpPoolData{}

	for k, v := range ths.m {
		data[k] = v
	}

	return data
}

func (ths *PoolMap) GetById(id string) (*models.DhcpPoolData, error) {
	ths.RLock()
	defer ths.RUnlock()

	if pool, ok := ths.m[id]; ok {
		return &pool, nil
	}

	return nil, fmt.Errorf("pool not found")
}

func (ths *PoolMap) GetByType(poolType int) (*models.DhcpPoolData, error) {
	ths.RLock()
	defer ths.RUnlock()

	for _, p := range ths.m {
		if p.Type == poolType {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("pool not found")
}

func (ths *PoolMap) GetByIp(ip string) (*models.DhcpPoolData, error) {
	ths.RLock()
	defer ths.RUnlock()

	for _, p := range ths.m {
		if gipv4.Ip2long(p.RangeStart) <= gipv4.Ip2long(ip) &&
			gipv4.Ip2long(p.RangeEnd) >= gipv4.Ip2long(ip) {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("pool not found")
}
