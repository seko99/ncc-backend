package simulator

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/events/interfaces"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	repository2 "code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"fmt"
	"time"
)

const (
	EventTypeSimulator = "SimulatorEvent"

	JobStateInProgress = "inProgress"
	JobStateDone       = "done"

	JobTypeCreateCustomers = "createCustomers"
	JobTypeCreateDevices   = "createDevices"
	JobTypeCreateMapNodes  = "createMapNodes"
	JobTypeCreateLeases    = "createLeases"
)

type BrasParams struct {
	SendInterims bool
}

type Simulator struct {
	log                  *zero.Logger
	storage              *psqlstorage.Storage
	events               interfaces.Events
	customerGroupsRepo   repository2.CustomerGroups
	citiesRepo           repository2.Cities
	streetsRepo          repository2.Streets
	customerRepo         repository2.Customers
	serviceInternetRepo  repository2.ServiceInternet
	paymentsRepo         repository2.Payments
	scoresRepo           repository2.Scores
	vendorsRepo          repository2.Vendors
	hardwareModelsRepo   repository2.HardwareModels
	paymentTypesRepo     repository2.PaymentTypes
	mapNodesRepo         repository2.MapNodes
	devicesRepo          repository2.Devices
	deviceStatesRepo     repository2.DeviceStates
	ifacesRepo           repository2.DeviceInterfaces
	ifaceStatesRepo      repository2.DeviceInterfaceStates
	dhcpBindingsRepo     repository2.DhcpBindings
	dhcpLeasesRepo       repository2.DhcpLeases
	sessionsRepo         repository2.Sessions
	nasesRepo            repository2.Nases
	nasTypesRepo         repository2.NasTypes
	feesRepo             repository2.Fees
	contractsRepo        repository2.Contracts
	paymentSystemsRepo   repository2.PaymentSystems
	usersRepo            repository2.Users
	radiusVendorsRepo    repository2.RadiusVendors
	radiusAttributesRepo repository2.RadiusAttributes
	issueTypesRepo       repository2.IssueTypes
	issueUrgenciesRepo   repository2.IssueUrgencies
	issuesRepo           repository2.Issues
	issueActionsRepo     repository2.IssueActions
	dhcpPoolsRepo        repository2.DhcpPools
	ipPoolRepo           repository2.IpPools
	lat                  float64
	lng                  float64

	sessionCache   repository2.Sessions
	leasesCache    repository2.DhcpLeases
	customersCache repository2.Customers
	nasesCache     repository2.Nases

	brasParams BrasParams
}

func NewSimulator(
	cfg *config.Config,
	log *zero.Logger,
	storage *psqlstorage.Storage,
	events interfaces.Events,
	customerGroupsRepo repository2.CustomerGroups,
	citiesRepo repository2.Cities,
	streetsRepo repository2.Streets,
	customerRepo repository2.Customers,
	serviceInternetRepo repository2.ServiceInternet,
	vendorsRepo repository2.Vendors,
	hardwareModelsRepo repository2.HardwareModels,
	paymentTypesRepo repository2.PaymentTypes,
	mapNodesRepo repository2.MapNodes,
	devicesRepo repository2.Devices,
	deviceStatesRepo repository2.DeviceStates,
	ifacesRepo repository2.DeviceInterfaces,
	ifaceStatesRepo repository2.DeviceInterfaceStates,
	dhcpBindingsRepo repository2.DhcpBindings,
	dhcpLeasesRepo repository2.DhcpLeases,
	sessionsRepo repository2.Sessions,
	nasesRepo repository2.Nases,
	nasTypesRepo repository2.NasTypes,
	paymentsRepo repository2.Payments,
	feesRepo repository2.Fees,
	contractsRepo repository2.Contracts,
	paymentSystemsRepo repository2.PaymentSystems,
	usersRepo repository2.Users,
	radiusVendorsRepo repository2.RadiusVendors,
	radiusAttributesRepo repository2.RadiusAttributes,
	dhcpPoolsRepo repository2.DhcpPools,
	ipPoolRepo repository2.IpPools,
	issueTypesRepo repository2.IssueTypes,
	issueUrgenciesRepo repository2.IssueUrgencies,
	issuesRepo repository2.Issues,
	issueActionsRepo repository2.IssueActions,

	sessionCache repository2.Sessions,
	leasesCache repository2.DhcpLeases,
	customersCache repository2.Customers,
	nasesCache repository2.Nases,
) *Simulator {
	simulator := &Simulator{
		log:                  log,
		storage:              storage,
		events:               events,
		customerGroupsRepo:   customerGroupsRepo,
		citiesRepo:           citiesRepo,
		streetsRepo:          streetsRepo,
		customerRepo:         customerRepo,
		serviceInternetRepo:  serviceInternetRepo,
		vendorsRepo:          vendorsRepo,
		hardwareModelsRepo:   hardwareModelsRepo,
		paymentTypesRepo:     paymentTypesRepo,
		mapNodesRepo:         mapNodesRepo,
		devicesRepo:          devicesRepo,
		deviceStatesRepo:     deviceStatesRepo,
		ifacesRepo:           ifacesRepo,
		ifaceStatesRepo:      ifaceStatesRepo,
		dhcpBindingsRepo:     dhcpBindingsRepo,
		dhcpLeasesRepo:       dhcpLeasesRepo,
		sessionsRepo:         sessionsRepo,
		nasesRepo:            nasesRepo,
		nasTypesRepo:         nasTypesRepo,
		paymentsRepo:         paymentsRepo,
		feesRepo:             feesRepo,
		contractsRepo:        contractsRepo,
		paymentSystemsRepo:   paymentSystemsRepo,
		usersRepo:            usersRepo,
		radiusVendorsRepo:    radiusVendorsRepo,
		radiusAttributesRepo: radiusAttributesRepo,
		issueTypesRepo:       issueTypesRepo,
		issueUrgenciesRepo:   issueUrgenciesRepo,
		issuesRepo:           issuesRepo,
		issueActionsRepo:     issueActionsRepo,
		dhcpPoolsRepo:        dhcpPoolsRepo,

		sessionCache:   sessionCache,
		leasesCache:    leasesCache,
		customersCache: customersCache,
		nasesCache:     nasesCache,
	}

	simulator.brasParams.SendInterims = true

	go simulator.radiusInterimUpdate(dto.RadiusUsecaseRequest{
		Secret:        cfg.Radius.Test.Secret,
		NasIP:         cfg.Radius.Test.Nas.Ip,
		NasIdentifier: cfg.Radius.Test.Nas.Identifier,
		Auth:          cfg.Radius.Test.Auth,
		Acct:          cfg.Radius.Test.Acct,
	})

	return simulator
}

func (ths *Simulator) SetBRASParams(params BrasParams) {
	ths.brasParams = params
}

func (ths *Simulator) GetSessionCache() ([]models.SessionData, error) {
	sessions, err := ths.sessionCache.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get session cache: %w", err)
	}
	return sessions, nil
}

func (ths *Simulator) GetLeasesCache() ([]models.LeaseData, error) {
	leases, err := ths.leasesCache.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get leases cache: %w", err)
	}
	return leases, nil
}

func (ths *Simulator) GetCustomerCache() ([]models.CustomerData, error) {
	customers, err := ths.customersCache.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get customers cache: %w", err)
	}
	return customers, nil
}

func (ths *Simulator) GetNASCache() ([]models.NasData, error) {
	nases, err := ths.nasesCache.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get nases cache: %w", err)
	}
	return nases, nil
}

func (ths *Simulator) ClearFakeData(req dto.FakeDataClearUsecaseRequest) error {

	if req.ClearGeo {
		err := ths.clearMapNodes()
		if err != nil {
			return fmt.Errorf("can't clear map nodes: %w", err)
		}
	}

	if req.ClearDevices {
		err := ths.clearDevices()
		if err != nil {
			return fmt.Errorf("can't clear devices: %w", err)
		}
	}

	if req.ClearCustomers {
		err := ths.clearCustomers()
		if err != nil {
			return fmt.Errorf("can't clear customers: %w", err)
		}
	}

	return nil
}

func (ths *Simulator) CreateFakeData(req dto.FakeDataCreateUsecaseRequest) error {

	if req.CreateCustomers {
		err := ths.createCustomers(req.MaxCustomers)
		if err != nil {
			return fmt.Errorf("can't create customers: %w", err)
		}
	}

	if req.CreateContracts {
		err := ths.createContracts()
		if err != nil {
			return fmt.Errorf("can't create contracts: %w", err)
		}
	}

	if req.CreateGeo {
		addresses, err := ths.getOsmData(req.LeftUpper, req.RightBottom, req.MinBuildLevel, req.MaxMapNodes)
		if err != nil {
			return fmt.Errorf("can't get OSM data: %w", err)
		}

		err = ths.createStreets(addresses)
		if err != nil {
			return fmt.Errorf("can't create streets: %w", err)
		}

		err = ths.createMapNodes(addresses)
		if err != nil {
			return fmt.Errorf("can't create map nodes: %w", err)
		}

		if req.CreateDevices {
			err = ths.createDevices()
			if err != nil {
				return fmt.Errorf("can't create devices: %w", err)
			}
			if req.DistributeCustomers {
				err := ths.distributeCustomers()
				if err != nil {
					return fmt.Errorf("can't distribute customers: %w", err)
				}

				err = ths.createBindings()
				if err != nil {
					return fmt.Errorf("can't create bindings: %w", err)
				}

				if req.CreateLeases {
					err := ths.createLeases()
					if err != nil {
						return fmt.Errorf("can't create leases: %w", err)
					}
				}
			}
		}

		if req.CreateSessions {
			err := ths.createSessions()
			if err != nil {
				return fmt.Errorf("can't create sessions: %w", err)
			}
		}

		err = ths.UpdateMap()
		if err != nil {
			ths.log.Error("Can't update map: %v", err)
		}
	}

	startTime := time.Now().Add(-30 * 24 * time.Hour)

	if req.CreatePayments {
		err := ths.createPayments(startTime)
		if err != nil {
			return fmt.Errorf("can't create payments: %w", err)
		}
	}

	if req.CreateFees {
		err := ths.createFees(startTime)
		if err != nil {
			return fmt.Errorf("can't create fees: %w", err)
		}
	}

	return nil
}
