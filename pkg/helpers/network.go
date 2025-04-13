package helpers

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"fmt"
	"github.com/gogf/gf/net/gipv4"
)

func InPool(pool models.IpPoolData, ip string) bool {
	ipLong := gipv4.Ip2long(ip)
	if ipLong >= gipv4.Ip2long(pool.PoolStart) && ipLong <= gipv4.Ip2long(pool.PoolEnd) {
		return true
	}
	return false
}

func GetPoolByIP(pools []models.IpPoolData, ip string) (models.IpPoolData, error) {
	for _, p := range pools {
		if InPool(p, ip) {
			return p, nil
		}
	}
	return models.IpPoolData{}, fmt.Errorf("pool not found for %s", ip)
}
