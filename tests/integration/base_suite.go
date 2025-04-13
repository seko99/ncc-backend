package integration

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/fixtures"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository/memory"
	psql2 "code.evixo.ru/ncc/ncc-backend/pkg/repository/psql"
	psqlstorage "code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"context"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	dblogger "gorm.io/gorm/logger"
	"strconv"
	"time"
)

const (
	RabbitDelay = 20 * time.Second
)

type BaseTestSuite struct {
	suite.Suite
	ctx    context.Context
	log    logger.Logger
	cfg    *config.Config
	events *events.Events

	eventsEnabled bool

	customersRepo        repository.Customers
	sessionsRepo         repository.Sessions
	sessionsLogRepo      repository.SessionsLog
	leasesRepo           repository.DhcpLeases
	poolsRepo            repository.DhcpPools
	bindingsRepo         repository.DhcpBindings
	nasesRepo            repository.Nases
	nasTypesRepo         repository.NasTypes
	radiusVendorsRepo    repository.RadiusVendors
	radiusAttributesRepo repository.RadiusAttributes

	sessionCache   repository.Sessions
	leasesCache    repository.DhcpLeases
	customersCache repository.Customers
	nasesCache     repository.Nases

	serviceInternetRepo repository.ServiceInternet

	informingsRepo              repository.Informings
	informingLogRepo            repository.InformingLog
	informingsTestCustomersRepo repository.InformingsTestCustomers
}

func (ths *BaseTestSuite) SetupSuite() {
	ths.ctx = context.Background()
	ths.log = zero.NewLogger(zerolog.DebugLevel)

	image := "postgres:14.3"
	req := testcontainers.ContainerRequest{
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_DB":       "postgres",
		},
		ExposedPorts: []string{"5432/tcp"},
		Image:        image,
		WaitingFor: wait.ForExec([]string{"pg_isready"}).
			WithPollInterval(2 * time.Second).
			WithExitCodeMatcher(func(exitCode int) bool {
				return exitCode == 0
			}),
	}

	container, err := testcontainers.GenericContainer(ths.ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	require.NoError(ths.T(), err)

	host, err := container.Host(ths.ctx)
	require.NoError(ths.T(), err)

	mappedPort, err := container.MappedPort(ths.ctx, "5432")
	require.NoError(ths.T(), err)

	port, err := strconv.Atoi(mappedPort.Port())
	require.NoError(ths.T(), err)

	ths.log.Info("DB: %s:%d", host, port)

	p, err := nat.NewPort("", "5672")
	require.NoError(ths.T(), err)

	cfg := &config.Config{
		Db: config.DbConfig{
			Host:               host,
			Port:               port,
			User:               "postgres",
			Password:           "password",
			Name:               "postgres",
			MaxLifeTime:        time.Hour,
			MaxOpenConnections: 5,
			MaxIdleConnections: 5,
		},
		Dhcp: config.DhcpConfig{
			Listen: "127.0.0.1:1067",
		},
		Radius: config.RadiusConfig{
			Auth: config.RadiusServerConfig{
				Listen: "0.0.0.0:1812",
			},
			Acct: config.RadiusServerConfig{
				Listen: "0.0.0.0:1813",
			},
			Secret: "someTestSecret",
			Update: time.Second * 10,
		},
	}

	if ths.eventsEnabled {
		rabbitImage := "rabbitmq:3"
		rabbitReq := testcontainers.ContainerRequest{
			Env: map[string]string{
				"RABBITMQ_DEFAULT_USER": "user",
				"RABBITMQ_DEFAULT_PASS": "password",
			},
			ExposedPorts: []string{string(p)},
			Image:        rabbitImage,
			WaitingFor:   wait.ForListeningPort(p).WithStartupTimeout(5 * time.Minute),
		}

		rabbitContainer, err := testcontainers.GenericContainer(ths.ctx,
			testcontainers.GenericContainerRequest{
				ContainerRequest: rabbitReq,
				Started:          true,
			},
		)
		require.NoError(ths.T(), err)

		rabbitHost, err := rabbitContainer.Host(ths.ctx)
		require.NoError(ths.T(), err)

		rabbitMappedPort, err := rabbitContainer.MappedPort(ths.ctx, p)
		require.NoError(ths.T(), err)

		rabbitPort, err := strconv.Atoi(rabbitMappedPort.Port())
		require.NoError(ths.T(), err)

		cfg.Queue = config.QueueConfig{
			Host:     rabbitHost,
			Port:     rabbitPort,
			User:     "user",
			Password: "password",
		}

		time.Sleep(RabbitDelay)
		ths.log.Info("Queue: %s:%d", rabbitHost, rabbitPort)

		ths.events, err = events.NewBroadcastEvents(cfg, ths.log, uuid.NewString(), events.BroadcastEvents)
		require.NoError(ths.T(), err)
		require.NotNil(ths.T(), ths.events)
	}

	ths.cfg = cfg

	storage := psqlstorage.NewStorage(cfg, ths.log, psqlstorage.WithLogLevel(dblogger.Error))
	err = storage.Connect()
	require.NoError(ths.T(), err)

	res := storage.GetDB().Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	require.NoError(ths.T(), res.Error)

	err = storage.Migrate()
	require.NoError(ths.T(), err)

	ths.customersRepo = psql2.NewCustomers(storage, ths.events)
	ths.nasesRepo = psql2.NewNases(storage, ths.events)
	ths.nasTypesRepo = psql2.NewNasTypes(storage)
	ths.radiusVendorsRepo = psql2.NewRadiusVendors(storage)
	ths.radiusAttributesRepo = psql2.NewRadiusAttributes(storage)
	ths.leasesRepo = psql2.NewDhcpLeases(storage, ths.events)
	ths.poolsRepo = psql2.NewDhcpPools(storage, ths.events)
	ths.bindingsRepo = psql2.NewDhcpBindings(storage, ths.events)
	ths.sessionsRepo = psql2.NewSessions(storage, ths.events)
	ths.sessionsLogRepo = psql2.NewSessionsLog(storage)
	ths.informingsRepo = psql2.NewInformings(storage)
	ths.informingLogRepo = psql2.NewInformingLog(storage)
	ths.informingsTestCustomersRepo = psql2.NewInformingsTestCustomers(storage)
	ths.serviceInternetRepo = psql2.NewServiceInternet(storage, ths.events)

	customers, err := ths.customersRepo.Get()
	assert.NoError(ths.T(), err)
	assert.Empty(ths.T(), customers)

	for _, c := range fixtures.FakeCustomers() {
		err := ths.customersRepo.Create(c)
		require.NoError(ths.T(), err)
	}

	for _, i := range fixtures.FakeInformings() {
		err := ths.informingsRepo.Create(i)
		require.NoError(ths.T(), err)
	}

	for _, v := range fixtures.FakeRadiusVendors() {
		err := ths.radiusVendorsRepo.Create(v)
		require.NoError(ths.T(), err)
	}

	for _, a := range fixtures.FakeRadiusAttrs() {
		err := ths.radiusAttributesRepo.Create(a)
		require.NoError(ths.T(), err)
	}

	for _, nt := range fixtures.FakeNasTypes() {
		err := ths.nasTypesRepo.Create(nt)
		require.NoError(ths.T(), err)
	}

	err = ths.poolsRepo.Create(fixtures.FakePools())
	require.NoError(ths.T(), err)

	err = ths.nasesRepo.Create(fixtures.FakeNases())
	require.NoError(ths.T(), err)

	for _, l := range fixtures.FakeLeases() {
		err := ths.leasesRepo.Create(l)
		require.NoError(ths.T(), err)
	}

	leases, err := ths.leasesRepo.Get()
	require.NoError(ths.T(), err)
	require.Len(ths.T(), leases, len(fixtures.FakeLeases()))

	err = ths.sessionsRepo.Create(fixtures.FakeSessions())
	require.NoError(ths.T(), err)

	sessionCache, err := memory.NewSessions(ths.log, ths.sessionsRepo, ths.events)
	require.NoError(ths.T(), err)
	customerCache, err := memory.NewCustomers(ths.log, ths.customersRepo, ths.events)
	require.NoError(ths.T(), err)
	leasesCache, err := memory.NewDhcpLeases(ths.log, ths.leasesRepo, ths.events)
	require.NoError(ths.T(), err)
	nasCache, err := memory.NewNases(ths.log, ths.nasesRepo, ths.events)
	require.NoError(ths.T(), err)

	ths.sessionCache = sessionCache
	ths.leasesCache = leasesCache
	ths.customersCache = customerCache
	ths.nasesCache = nasCache
}
