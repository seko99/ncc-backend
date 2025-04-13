package simulator

import (
	fixtures2 "code.evixo.ru/ncc/ncc-backend/pkg/fixtures"
	"fmt"
)

func (ths *Simulator) InitDictionaries() error {

	err := ths.initAddresses()
	if err != nil {
		return fmt.Errorf("can't init addresses: %w", err)
	}

	err = ths.initCustomerGroups()
	if err != nil {
		return fmt.Errorf("can't init customer groups: %w", err)
	}

	err = ths.initServiceInternet()
	if err != nil {
		return fmt.Errorf("can't init service Internet: %w", err)
	}

	err = ths.initVendors()
	if err != nil {
		return fmt.Errorf("can't init vendors: %w", err)
	}

	err = ths.initHardwareModels()
	if err != nil {
		return fmt.Errorf("can't init hardware models: %w", err)
	}

	err = ths.initPaymentTypes()
	if err != nil {
		return fmt.Errorf("can't init payment types: %w", err)
	}

	err = ths.initNasTypes()
	if err != nil {
		return fmt.Errorf("can't init NAS types: %w", err)
	}

	err = ths.initNases()
	if err != nil {
		return fmt.Errorf("can't init nases: %w", err)
	}

	/*	err = ths.initPaymentSystems()
		if err != nil {
			return fmt.Errorf("can't init payment systems: %w", err)
		}
	*/
	err = ths.initRadiusVendors()
	if err != nil {
		return fmt.Errorf("can't init RADIUS vendors: %w", err)
	}

	err = ths.initIssueTypes()
	if err != nil {
		return fmt.Errorf("can't init issue types: %w", err)
	}

	err = ths.initIssueUrgencies()
	if err != nil {
		return fmt.Errorf("can't init issue urgencies: %w", err)
	}

	err = ths.initDhcpPools()
	if err != nil {
		return fmt.Errorf("can't init DHCP pools: %w", err)
	}

	return nil
}

func (ths *Simulator) initDhcpPools() error {
	ths.log.Info("Creating DHCP pools...")

	count := len(fixtures2.FakeDhcpPools())
	for _, p := range fixtures2.FakeDhcpPools() {
		err := ths.dhcpPoolsRepo.Upsert(p)
		if err != nil {
			ths.log.Error("Can't create DHCP pool: %v", err)
			continue
		}
		count--
	}

	if count > 0 {
		return fmt.Errorf("can't create some of DHCP pools")
	}

	return nil
}

func (ths *Simulator) initPaymentSystems() error {
	ths.log.Info("Creating payment systems...")

	customers, err := ths.customerRepo.Get(1)
	if err != nil {
		return fmt.Errorf("can't get customers: %w", err)
	}

	users, err := ths.usersRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get users: %w", err)
	}

	customerId := customers[0].Id
	userId := users[0].Id

	for _, p := range fixtures2.FakePaymentSystems(userId, customerId) {
		err := ths.paymentSystemsRepo.Upsert(p)
		if err != nil {
			ths.log.Error("Can't create payment system: %v", err)
		}
	}
	return nil
}

func (ths *Simulator) initNasTypes() error {
	ths.log.Info("Creating NAS types...")

	for _, v := range fixtures2.FakeRadiusVendors() {
		err := ths.radiusVendorsRepo.Create(v)
		if err != nil {
			return fmt.Errorf("can't create RADIUS vendor: %w", err)
		}
	}

	for _, a := range fixtures2.FakeRadiusAttrs() {
		err := ths.radiusAttributesRepo.Create(a)
		if err != nil {
			return fmt.Errorf("can't create RADIUS attribute: %w", err)
		}
	}

	created := 0
	for _, t := range fixtures2.FakeNasTypes() {
		err := ths.nasTypesRepo.Upsert(t)
		if err != nil {
			ths.log.Error("Can't create NAS type: %v", err)
			continue
		}
		created++
	}

	if created == 0 {
		return fmt.Errorf("can't create at least one NAS type")
	}

	return nil
}

func (ths *Simulator) initIssueTypes() error {
	ths.log.Info("Creating issue types...")

	count := len(fixtures2.FakeIssueTypes())
	for _, t := range fixtures2.FakeIssueTypes() {
		err := ths.issueTypesRepo.Upsert(t)
		if err != nil {
			ths.log.Error("Can't create issue type: %v", err)
			continue
		}
		count--
	}

	if count > 0 {
		return fmt.Errorf("can't create some of issue types")
	}

	return nil
}

func (ths *Simulator) initIssueUrgencies() error {
	ths.log.Info("Creating issue urgencies...")

	count := len(fixtures2.FakeIssueUrgencies())
	for _, u := range fixtures2.FakeIssueUrgencies() {
		err := ths.issueUrgenciesRepo.Upsert(u)
		if err != nil {
			ths.log.Error("Can't create issue urgency: %v", err)
			continue
		}
		count--
	}

	if count > 0 {
		return fmt.Errorf("can't create some of issue urgencies")
	}

	return nil
}

func (ths *Simulator) initRadiusVendors() error {
	ths.log.Info("Creating RADIUS vendors...")

	count := len(fixtures2.FakeRadiusVendors())
	for _, v := range fixtures2.FakeRadiusVendors() {
		err := ths.radiusVendorsRepo.Upsert(v)
		if err != nil {
			ths.log.Error("Can't create RADIUS vendor: %v", err)
			continue
		}
		count--
	}

	if count > 0 {
		return fmt.Errorf("can't create some RADIUS vendors")
	}

	count = len(fixtures2.FakeRadiusAttrs())
	for _, a := range fixtures2.FakeRadiusAttrs() {
		err := ths.radiusAttributesRepo.Upsert(a)
		if err != nil {
			ths.log.Error("Can't create RADIUS attribute: %v", err)
			continue
		}
		count--
	}

	if count > 0 {
		return fmt.Errorf("can't create some RADIUS attributes")
	}

	return nil
}

func (ths *Simulator) initNases() error {
	ths.log.Info("Creating NASes...")

	created := 0
	for _, n := range fixtures2.FakeNases() {
		err := ths.nasesRepo.Upsert(n)
		if err != nil {
			ths.log.Error("Can't create NAS: %v", err)
			continue
		}
		created++
	}

	if created == 0 {
		return fmt.Errorf("can't create at least one NAS")
	}

	return nil
}

func (ths *Simulator) initPaymentTypes() error {
	ths.log.Info("Creating payment types...")

	for _, ptype := range fixtures2.FakePaymentTypes() {
		err := ths.paymentTypesRepo.Upsert(ptype)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ths *Simulator) initCustomerGroups() error {
	ths.log.Info("Creating customer groups...")

	for _, group := range fixtures2.FakeCustomerGroups() {
		err := ths.customerGroupsRepo.Upsert(group)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ths *Simulator) initAddresses() error {
	ths.log.Info("Creating addresses...")

	for _, city := range fixtures2.FakeCities() {
		err := ths.citiesRepo.Upsert(city)
		if err != nil {
			return err
		}
	}

	for _, street := range fixtures2.FakeStreets() {
		err := ths.streetsRepo.Upsert(street)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ths *Simulator) initHardwareModels() error {
	ths.log.Info("Creating hardware models...")

	for _, h := range fixtures2.FakeHardwareModels() {
		err := ths.hardwareModelsRepo.Upsert(h)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ths *Simulator) initVendors() error {
	ths.log.Info("Creating vendors...")

	for _, v := range fixtures2.FakeVendors() {
		err := ths.vendorsRepo.Upsert(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ths *Simulator) initServiceInternet() error {
	ths.log.Info("Creating services...")

	for _, service := range fixtures2.FakeServicesInternet() {
		err := ths.serviceInternetRepo.Upsert(service)
		if err != nil {
			return err
		}
	}
	return nil
}
