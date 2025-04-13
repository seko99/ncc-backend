package simulator

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"fmt"
	"github.com/bxcodec/faker/v4"
	"github.com/gogf/gf/net/gipv4"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

const (
	DevicePortCount = 8
)

func (ths *Simulator) clearIfaces() error {
	ths.log.Info("Clearing ifaces...")

	r := ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_iface")
	if r.Error != nil {
		return fmt.Errorf("can't delete ifaces: %w", r.Error)
	}

	r = ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_iface_state")
	if r.Error != nil {
		return fmt.Errorf("can't delete iface states: %w", r.Error)
	}

	return nil
}

func (ths *Simulator) clearPONONU() error {
	ths.log.Info("Clearing PON ONU...")

	r := ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_pon_onu_state")
	if r.Error != nil {
		return fmt.Errorf("can't delete PON ONU states: %w", r.Error)
	}

	r = ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_pon_onu")
	if r.Error != nil {
		return fmt.Errorf("can't delete PON ONUs: %w", r.Error)
	}

	return nil
}

func (ths *Simulator) clearFDB() error {
	ths.log.Info("Clearing FDB...")

	r := ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_fdb")
	if r.Error != nil {
		return fmt.Errorf("can't delete FDB: %w", r.Error)
	}

	return nil
}

func (ths *Simulator) clearDevices() error {
	ths.log.Info("Clearing devices...")

	err := ths.clearLeases()
	if err != nil {
		return fmt.Errorf("can't clear leases: %w", err)
	}

	err = ths.clearBindings()
	if err != nil {
		return fmt.Errorf("can't clear bindings: %w", err)
	}

	err = ths.clearIfaces()
	if err != nil {
		return fmt.Errorf("can't clear ifaces: %w", err)
	}

	err = ths.clearFDB()
	if err != nil {
		return fmt.Errorf("can't clear FDB: %w", err)
	}

	err = ths.clearPONONU()
	if err != nil {
		return fmt.Errorf("can't clear PON ONU: %w", err)
	}

	r := ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_device")
	if r.Error != nil {
		return fmt.Errorf("can't delete devices: %w", err)
	}

	r = ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_device_state")
	if r.Error != nil {
		return fmt.Errorf("can't delete device states: %w", err)
	}

	return nil
}

func (ths *Simulator) createDevices() error {

	mapNodes, err := ths.mapNodesRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get map nodes: %w", err)
	}

	ths.log.Info("Creating %d devices...", len(mapNodes))

	hardwareModels, err := ths.hardwareModelsRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get hardware models: %w", err)
	}

	cities, err := ths.citiesRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get cities: %w", err)
	}

	ip := gipv4.Ip2long("10.0.30.10")
	svid := 1000

	for nodeIdx, m := range mapNodes {
		deviceId := uuid.NewString()
		cityId := cities[0].Id
		modelId := hardwareModels[rand.Intn(len(hardwareModels))].Id
		stateId := uuid.NewString()

		err := ths.deviceStatesRepo.Create(models2.DeviceStateData{
			CommonData: models2.CommonData{
				Id: stateId,
			},
			StatusSnmp:    models2.DeviceStatusOK,
			StatusIcmp:    models2.DeviceStatusOK,
			StatusUpdated: time.Now(),
		})
		if err != nil {
			ths.log.Error("Can't create device state: %v", err)
			continue
		}

		err = ths.devicesRepo.Create(models2.DeviceData{
			CommonData: models2.CommonData{
				Id: deviceId,
			},
			MapNodeId:     models2.NewNullUUID(m.Id),
			ModelId:       models2.NewNullUUID(modelId),
			CityId:        models2.NewNullUUID(cityId),
			StreetId:      m.StreetId,
			Build:         m.Build,
			Lat:           m.Lat,
			Lng:           m.Lng,
			PortCount:     DevicePortCount,
			Svid:          svid,
			Ip:            gipv4.Long2ip(ip),
			DeviceStateId: models2.NewNullUUID(stateId),
			Remote:        faker.MacAddress(),
		})
		if err != nil {
			ths.log.Error("Can't create device at %s, %s: %v", m.Street.Name, m.Build, err)
		}

		for i := 0; i < DevicePortCount; i++ {
			port := i + 1

			ths.log.Info("Creating interface %d/%d at node %d/%d...", port, DevicePortCount, nodeIdx+1, len(mapNodes))

			ifaceStateId := uuid.NewString()

			status := models2.DeviceIfaceStatusDown
			speed := uint64(0)
			descr := fmt.Sprintf("port %d", port)

			if port == DevicePortCount {
				status = models2.DeviceIfaceStatusUp
				speed = models2.DeviceIfaceSpeed1000
				descr = "UPLINK"
			}

			err := ths.ifaceStatesRepo.Create(models2.IfaceStateData{
				CommonData: models2.CommonData{
					Id: ifaceStateId,
				},
				LastStatus: models2.DeviceIfaceStatusDown,
				LastChange: time.Now(),
				OperStatus: status,
				Speed:      speed,
			})
			if err != nil {
				ths.log.Error("Can't create interface state: %v", err)
				continue
			}

			err = ths.ifacesRepo.Create(models2.IfaceData{
				DeviceId:     models2.NewNullUUID(deviceId),
				IfaceStateId: models2.NewNullUUID(ifaceStateId),
				Iface:        fmt.Sprintf("%d", port),
				Descr:        descr,
				Port:         port,
				Type:         models2.DeviceIfaceTypeNormal,
			})
			if err != nil {
				ths.log.Error("Can't create interface: %v", err)
			}
		}

		ip++
		svid++

		err = ths.events.PublishEvent(events.Event{
			Type: EventTypeSimulator,
			Payload: map[string]interface{}{
				"job_type":  JobTypeCreateDevices,
				"job_state": JobStateInProgress,
				"max":       len(mapNodes),
				"progress":  nodeIdx,
			},
		})
		if err != nil {
			ths.log.Error("Can't publish event: %v", err)
		}
	}

	err = ths.events.PublishEvent(events.Event{
		Type: EventTypeSimulator,
		Payload: map[string]interface{}{
			"job_type":  JobTypeCreateDevices,
			"job_state": JobStateDone,
			"max":       len(mapNodes),
			"progress":  len(mapNodes),
		},
	})
	if err != nil {
		ths.log.Error("Can't publish event: %v", err)
	}

	return nil
}
