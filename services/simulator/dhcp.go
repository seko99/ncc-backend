package simulator

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"fmt"
	"github.com/bxcodec/faker/v4"
	"github.com/gogf/gf/net/gipv4"
	"time"
)

func (ths *Simulator) UpdateLeases() error {
	leases, err := ths.dhcpLeasesRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get leases: %w", err)
	}

	ths.log.Info("Updating %d leases...", len(leases))

	for _, l := range leases {
		l.LeaseTime = 360
		l.Expire = time.Now().Add(180 * time.Second)
		l.Status = models2.DhcpLeaseStatusAccepted

		err := ths.dhcpLeasesRepo.Update(l)
		if err != nil {
			ths.log.Error("Can't update lease: %v", err)
		}
	}

	return nil
}

func (ths *Simulator) clearLeases() error {
	ths.log.Info("Clearing leases...")

	err := ths.dhcpLeasesRepo.DeleteAll()
	if err != nil {
		return fmt.Errorf("can't delete leases: %w", err)
	}

	return nil
}

func (ths *Simulator) createLeases() error {

	bindings, err := ths.dhcpBindingsRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get bindings: %w", err)
	}

	ths.log.Info("Creating %d leases...", len(bindings))

	ip := gipv4.Ip2long("172.45.0.10")

	for _, b := range bindings {
		err := ths.dhcpLeasesRepo.Create(models2.LeaseData{
			CustomerId: b.CustomerId,
			DeviceId:   b.DeviceId,
			Mac:        b.Mac,
			Cvid:       b.Cvid,
			Port:       b.Port,
			Remote:     b.Remote,
			Start:      time.Now(),
			Expire:     time.Now().Add(3 * time.Minute),
			Ip:         gipv4.Long2ip(ip),
			Subnet:     "255.255.0.0",
			Router:     "172.45.0.1",
			Dns1:       "8.8.8.8",
			Dns2:       "1.1.1.1",
			Status:     models2.DhcpLeaseStatusAccepted,
		})
		if err != nil {
			ths.log.Error("Can't create lease: %v", err)
			continue
		}

		ip++

		iface, err := ths.ifacesRepo.GetByDeviceIdAndPort(b.DeviceId.UUID.String(), b.Port)
		if err != nil {
			ths.log.Error("Can't find interface: %v", err)
			continue
		}

		err = ths.ifaceStatesRepo.Update(models2.IfaceStateData{
			CommonData: models2.CommonData{
				Id: iface.IfaceStateId.UUID.String(),
			},
			OperStatus: models2.DeviceIfaceStatusUp,
			Speed:      models2.DeviceIfaceSpeed100,
			LastChange: time.Now(),
			LastStatus: models2.DeviceIfaceStatusDown,
		})
		if err != nil {
			ths.log.Error("Can't update interface state: %v", err)
		}
	}

	return nil
}

func (ths *Simulator) clearBindings() error {
	ths.log.Info("Clearing bindings...")

	r := ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_dhcp_binding")
	if r.Error != nil {
		return fmt.Errorf("can't delete bindings: %w", r.Error)
	}

	return nil
}

func (ths *Simulator) createBindings() error {

	customers, err := ths.customerRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get customers: %w", err)
	}

	devices, err := ths.devicesRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get devices: %w", err)
	}

	devicePorts := map[string]int{}

	ths.log.Info("Creating %d bindings...", len(customers))

	for _, c := range customers {
		for _, d := range devices {
			if d.StreetId == c.StreetId && d.Build == c.Build {
				port, ok := devicePorts[d.Id]
				if !ok {
					devicePorts[d.Id] = 1
					port = 1
				}
				err := ths.dhcpBindingsRepo.Create([]models2.DhcpBindingData{
					{
						CommonData: models2.CommonData{},
						CustomerId: models2.NewNullUUID(c.Id),
						DeviceId:   models2.NewNullUUID(d.Id),
						Remote:     d.Remote,
						Mac:        faker.MacAddress(),
						Port:       port,
						Cvid:       d.Svid,
					},
				})
				if err != nil {
					ths.log.Error("Can't create binding: %v", err)
				}
				devicePorts[d.Id]++
				break
			}
		}
	}

	return nil
}
