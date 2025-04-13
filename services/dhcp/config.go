package dhcp

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"encoding/json"
	"fmt"
	"time"
)

func (ths *Server) configUpdater() {
	for {
		err := ths.events.PublishRequest(events.Event{
			Type: ConfigRequest,
		}, ConfigResponse, func(event events.Event, params ...interface{}) {
			ths.updateConfigs(event.Payload)
		})
		if err != nil {
			ths.log.Error("can't get config: %v", err)
		}

		time.Sleep(ths.cfg.Radius.Update)
	}
}

func (ths *Server) getLeases(data interface{}) (int, error) {
	b, _ := json.Marshal(data)

	m := map[string]models.LeaseData{}
	err := json.Unmarshal(b, &m)
	if err != nil {
		return 0, fmt.Errorf("can't unmarshal lease map: %w", err)
	}
	ths.serverLeases.Set(m)

	return len(m), nil
}

func (ths *Server) getBindings(data interface{}) (int, error) {
	m := map[string]models.DhcpBindingData{}
	b, _ := json.Marshal(data)
	err := json.Unmarshal(b, &m)
	if err != nil {
		return 0, fmt.Errorf("can't unmarshal bindings map: %w", err)
	}
	ths.bindingMap.Set(m)

	return len(m), nil
}

func (ths *Server) getPools(data interface{}) (int, error) {
	m := map[string]models.DhcpPoolData{}
	b, _ := json.Marshal(data)
	err := json.Unmarshal(b, &m)
	if err != nil {
		return 0, fmt.Errorf("can't unmarshal pool map: %w", err)
	}
	ths.poolMap.Set(m)

	return len(m), nil
}

func (ths *Server) updateConfigs(payload map[string]interface{}) {
	mapLen, err := ths.getPools(payload["pools"])
	if err != nil {
		ths.log.Error("Can't get pools: %v", err)
	} else {
		ths.log.Info("Received pool map: %d", mapLen)
	}

	if mapLen > 0 && ths.startup {
		if mapLen > 0 {
			ths.startup = false
			ths.Wg.Done()
		}
	}

	mapLen, err = ths.getLeases(payload["leases"])
	if err != nil {
		ths.log.Error("Can't get leases: %v", err)
	} else {
		ths.log.Info("Received lease map: %d", mapLen)
	}

	mapLen, err = ths.getBindings(payload["bindings"])
	if err != nil {
		ths.log.Error("Can't get bindings: %v", err)
	} else {
		ths.log.Info("Received bindings map: %d", mapLen)
	}
}
