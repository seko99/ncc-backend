package integration

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/fixtures"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/services/radius"
	"context"
	"fmt"
	"github.com/gogf/gf/net/gipv4"
	"github.com/google/uuid"
	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	rad "layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"
	"layeh.com/radius/rfc2869"
	"net"
	"sync"
	"testing"
	"time"
)

const (
	LoginStormCount      = 1500
	CustomersCreateDelay = 60 * time.Second
	SessionsCreateDelay  = 60 * time.Second
)

type RadiusServerTestSuite struct {
	BaseTestSuite
	acctMap            map[string]models.SessionData
	radiusEventHandler *radius.EventHandler
	radiusServer       *radius.RadiusServer
	rsEvents           *events.Events
}

func (ths *RadiusServerTestSuite) TestRadius_00_SessionStart() {
	err := ths.sessionsRepo.DeleteAll()
	require.NoError(ths.T(), err)
	err = ths.sessionCache.DeleteAll()
	require.NoError(ths.T(), err)

	ths.radiusEventHandler, err = radius.NewRadiusEventHandler(
		ths.cfg,
		ths.log,
		ths.customersRepo,
		ths.nasesRepo,
		ths.leasesRepo,
		ths.sessionsRepo,
		ths.sessionsLogRepo,
		ths.serviceInternetRepo,
	)
	require.NoError(ths.T(), err)

	go ths.radiusEventHandler.Start()

	ths.radiusEventHandler.Wg.Wait()
	ths.log.Info("RadiusEventHandler ready")

	ths.rsEvents, err = events.NewEvents(ths.cfg, ths.log, uuid.NewString(), radius.Queue)
	require.NoError(ths.T(), err)

	ths.radiusServer = radius.NewRadiusServer(ths.cfg, ths.log, ths.rsEvents)

	go func() {
		err := ths.radiusServer.Start()
		require.NoError(ths.T(), err)
	}()

	ths.radiusServer.Wg.Wait()
	ths.log.Info("RadiusServer ready")

	nases, err := ths.nasesRepo.Get()
	nas := nases[0]

	customers, err := ths.customersRepo.Get()
	require.NoError(ths.T(), err)

	shouldReject := map[string]struct{}{
		"user2":  {}, // disabled
		"user17": {}, // wrong password
		"user18": {}, // disabled
		"user19": {}, // empty password
	}

	substPassword := map[string]string{
		"user17": "someWrongPassword",
		"user19": "",
	}

	var nasPort uint32 = 0
	var ip uint32 = gipv4.Ip2long("10.27.0.5")

	ths.acctMap = map[string]models.SessionData{}

	for _, c := range customers {
		password := c.Password

		if p, ok := substPassword[c.Login]; ok {
			password = p
		}

		packet := rad.New(rad.CodeAccessRequest, []byte(nas.Secret))
		err := rfc2865.UserName_SetString(packet, c.Login)
		require.NoError(ths.T(), err)
		err = rfc2865.UserPassword_SetString(packet, password)
		require.NoError(ths.T(), err)
		err = rfc2865.NASIPAddress_Set(packet, net.ParseIP(nas.Ip))
		require.NoError(ths.T(), err)
		err = rfc2865.NASIdentifier_SetString(packet, "local")
		require.NoError(ths.T(), err)
		err = rfc2865.NASPortType_Set(packet, rfc2865.NASPortType_Value_Ethernet)
		require.NoError(ths.T(), err)
		radResponse, err := rad.Exchange(context.Background(), packet, "127.0.0.1:1812")
		assert.NoError(ths.T(), err)
		assert.NotNil(ths.T(), radResponse)

		if _, ok := shouldReject[c.Login]; ok {
			assert.Equal(ths.T(), rad.CodeAccessReject, radResponse.Code)
		} else {
			assert.Equal(ths.T(), rad.CodeAccessAccept, radResponse.Code)
		}

		if radResponse.Code == rad.CodeAccessAccept {
			sessionId := uuid.NewString()

			data := models.SessionData{
				StartTime:     time.Now(),
				LastAlive:     time.Now(),
				AcctSessionId: sessionId,
				CustomerId:    models.NewNullUUID(c.Id),
				Login:         c.Login,
				Ip:            gipv4.Long2ip(ip),
				Nas:           nas,
				NasId:         models.NewNullUUID(nas.Id),
				NasPort:       nasPort,
			}

			packet := rad.New(rad.CodeAccountingRequest, []byte(nas.Secret))
			err = rfc2866.AcctStatusType_Set(packet, rfc2866.AcctStatusType_Value_Start)
			require.NoError(ths.T(), err)
			err = rfc2865.ServiceType_Add(packet, rfc2865.ServiceType_Value_FramedUser)
			require.NoError(ths.T(), err)
			err = rfc2865.FramedProtocol_Add(packet, rfc2865.FramedProtocol_Value_PPP)
			require.NoError(ths.T(), err)
			err = rfc2865.FramedIPAddress_Add(packet, net.ParseIP(gipv4.Long2ip(ip)))
			require.NoError(ths.T(), err)
			err = rfc2865.UserName_SetString(packet, c.Login)
			require.NoError(ths.T(), err)
			err = rfc2865.NASIPAddress_Set(packet, net.ParseIP(nas.Ip))
			require.NoError(ths.T(), err)
			err = rfc2865.NASIdentifier_AddString(packet, "local")
			require.NoError(ths.T(), err)
			err = rfc2865.NASPortType_Add(packet, rfc2865.NASPortType_Value_Ethernet)
			require.NoError(ths.T(), err)
			err = rfc2865.NASPort_Add(packet, rfc2865.NASPort(data.NasPort))
			require.NoError(ths.T(), err)

			err = rfc2866.AcctAuthentic_Add(packet, rfc2866.AcctAuthentic_Value_RADIUS)
			require.NoError(ths.T(), err)
			err = rfc2866.AcctSessionID_Add(packet, []byte(sessionId))
			require.NoError(ths.T(), err)
			err = rfc2866.AcctSessionTime_Set(packet, 0)
			require.NoError(ths.T(), err)

			radResponse, err = rad.Exchange(context.Background(), packet, "127.0.0.1:1813")
			assert.NoError(ths.T(), err)
			assert.NotNil(ths.T(), radResponse)

			ths.acctMap[c.Login] = data
			ip++
		}
	}

	time.Sleep(1 * time.Second)

	sessions, err := ths.sessionsRepo.Get()
	assert.NoError(ths.T(), err)

	for _, a := range ths.acctMap {
		found := false
		for _, s := range sessions {
			if s.Login == a.Login {
				found = true
				break
			}
		}
		assert.Equal(ths.T(), true, found, "Session for %s not exists", a.Login)
	}

	assert.Equal(ths.T(), len(customers)-len(shouldReject), len(sessions))
}

func (ths *RadiusServerTestSuite) TestRadius_01_SessionUpdate() {

	wg := sync.WaitGroup{}

	for _, session := range ths.acctMap {
		wg.Add(1)
		session.OctetsIn = 1024
		session.OctetsOut = 2048
		session.Duration = 10

		packet := rad.New(rad.CodeAccountingRequest, []byte(session.Nas.Secret))
		err := rfc2866.AcctStatusType_Set(packet, rfc2866.AcctStatusType_Value_InterimUpdate)
		require.NoError(ths.T(), err)
		err = rfc2865.ServiceType_Add(packet, rfc2865.ServiceType_Value_FramedUser)
		require.NoError(ths.T(), err)
		err = rfc2865.FramedProtocol_Add(packet, rfc2865.FramedProtocol_Value_PPP)
		require.NoError(ths.T(), err)
		err = rfc2865.FramedIPAddress_Add(packet, net.ParseIP(session.Ip))
		require.NoError(ths.T(), err)
		err = rfc2865.UserName_SetString(packet, session.Login)
		require.NoError(ths.T(), err)
		err = rfc2865.NASIPAddress_Set(packet, net.ParseIP(session.Nas.Ip))
		require.NoError(ths.T(), err)
		err = rfc2865.NASIdentifier_AddString(packet, "local")
		require.NoError(ths.T(), err)
		err = rfc2865.NASPortType_Add(packet, rfc2865.NASPortType_Value_Ethernet)
		require.NoError(ths.T(), err)
		err = rfc2865.NASPort_Add(packet, rfc2865.NASPort(session.NasPort))
		require.NoError(ths.T(), err)

		err = rfc2866.AcctAuthentic_Add(packet, rfc2866.AcctAuthentic_Value_RADIUS)
		require.NoError(ths.T(), err)
		err = rfc2866.AcctSessionID_Add(packet, []byte(session.AcctSessionId))
		require.NoError(ths.T(), err)
		err = rfc2866.AcctSessionTime_Set(packet, rfc2866.AcctSessionTime(session.Duration))
		require.NoError(ths.T(), err)
		err = rfc2866.AcctInputOctets_Set(packet, rfc2866.AcctInputOctets(session.OctetsIn))
		require.NoError(ths.T(), err)
		err = rfc2866.AcctOutputOctets_Set(packet, rfc2866.AcctOutputOctets(session.OctetsOut))
		require.NoError(ths.T(), err)

		radResponse, err := rad.Exchange(context.Background(), packet, "127.0.0.1:1813")
		assert.NoError(ths.T(), err)
		assert.NotNil(ths.T(), radResponse)

		ths.acctMap[session.Login] = session
		wg.Done()
	}

	wg.Wait()

	time.Sleep(1 * time.Second)

	sessions, err := ths.radiusEventHandler.GetSessionCache().Get()
	assert.NoError(ths.T(), err)

	for _, a := range ths.acctMap {
		found := false
		for _, s := range sessions {
			if s.Login == a.Login {
				found = true
				assert.Equal(ths.T(), a.OctetsIn, s.OctetsIn, "Login: %s", s.Login)
				assert.Equal(ths.T(), a.OctetsOut, s.OctetsOut, "Login: %s", s.Login)
				assert.Equal(ths.T(), a.Duration, s.Duration, "Login: %s", s.Login)
				break
			}
		}
		assert.Equal(ths.T(), true, found, "Session for %s not exists", a.Login)
	}
}

func (ths *RadiusServerTestSuite) TestRadius_02_LoginStorm() {
	shouldReject := map[string]struct{}{
		"user2":  {}, // disabled
		"user18": {}, // disabled
	}

	handler, err := radius.NewRadiusEventHandler(
		ths.cfg,
		ths.log,
		ths.customersRepo,
		ths.nasesRepo,
		ths.leasesRepo,
		ths.sessionsRepo,
		ths.sessionsLogRepo,
		ths.serviceInternetRepo,
	)
	require.NoError(ths.T(), err)

	ths.radiusEventHandler = handler

	go ths.radiusEventHandler.Start()

	ths.radiusEventHandler.Wg.Wait()
	ths.log.Info("RadiusEventHandler ready")

	nases, err := ths.nasesRepo.Get()
	nas := nases[0]

	var nasPort uint32 = 0
	var ip uint32 = gipv4.Ip2long("10.27.0.100")

	ths.acctMap = map[string]models.SessionData{}

	for i := 100; i < LoginStormCount; i++ {
		login := fmt.Sprintf("user%00d", i)

		err := ths.customersRepo.Create(models.CustomerData{
			CommonData: models.CommonData{
				Id:        uuid.NewString(),
				CreateTs:  time.Now(),
				CreatedBy: "16558104-13cb-4ffe-8dfe-2d32d1fc3acb",
			},
			Uid:                  random.New().String(12, "0123456789"),
			Login:                login,
			Password:             "testPassw0rd",
			Phone:                "70000000019",
			Deposit:              42,
			Credit:               0.0,
			BlockingState:        models.CustomerStateActive,
			ServiceInternetState: models.ServiceStateEnabled,
			ServiceInternet:      fixtures.Internet1,
			ServiceInternetId:    models.NewNullUUID(fixtures.Internet1.CommonData.Id),
			Flags: []models.CustomerFlagData{
				{
					Name: models.FieldSent,
					Val:  models.ExprTrue,
				},
			},
		})

		require.NoError(ths.T(), err)
	}

	ths.log.Info("Waiting for customers created...")
	time.Sleep(CustomersCreateDelay)

	customers, err := ths.customersRepo.Get()
	require.NoError(ths.T(), err)

	for _, c := range customers {
		password := c.Password

		packet := rad.New(rad.CodeAccessRequest, []byte(nas.Secret))
		err = rfc2865.UserName_SetString(packet, c.Login)
		require.NoError(ths.T(), err)
		err = rfc2865.UserPassword_SetString(packet, password)
		require.NoError(ths.T(), err)
		err = rfc2865.NASIPAddress_Set(packet, net.ParseIP(nas.Ip))
		require.NoError(ths.T(), err)
		err = rfc2865.NASIdentifier_SetString(packet, "local")
		require.NoError(ths.T(), err)
		err = rfc2865.NASPortType_Set(packet, rfc2865.NASPortType_Value_Ethernet)
		require.NoError(ths.T(), err)
		radResponse, err := rad.Exchange(context.Background(), packet, "127.0.0.1:1812")
		require.NoError(ths.T(), err)
		require.NotNil(ths.T(), radResponse)

		if _, ok := shouldReject[c.Login]; ok {
			require.Equal(ths.T(), rad.CodeAccessReject, radResponse.Code)
		} else {
			require.Equal(ths.T(), rad.CodeAccessAccept, radResponse.Code)
		}

		if radResponse.Code == rad.CodeAccessAccept {
			sessionId := uuid.NewString()

			data := models.SessionData{
				StartTime:     time.Now(),
				LastAlive:     time.Now(),
				AcctSessionId: sessionId,
				CustomerId:    models.NewNullUUID(c.Id),
				Login:         c.Login,
				Ip:            gipv4.Long2ip(ip),
				Nas:           nas,
				NasId:         models.NewNullUUID(nas.Id),
				NasPort:       nasPort,
			}

			packet := rad.New(rad.CodeAccountingRequest, []byte(nas.Secret))
			err = rfc2866.AcctStatusType_Set(packet, rfc2866.AcctStatusType_Value_Start)
			require.NoError(ths.T(), err)
			err = rfc2865.ServiceType_Add(packet, rfc2865.ServiceType_Value_FramedUser)
			require.NoError(ths.T(), err)
			err = rfc2865.FramedProtocol_Add(packet, rfc2865.FramedProtocol_Value_PPP)
			require.NoError(ths.T(), err)
			err = rfc2865.FramedIPAddress_Add(packet, net.ParseIP(gipv4.Long2ip(ip)))
			require.NoError(ths.T(), err)
			err = rfc2865.UserName_SetString(packet, c.Login)
			require.NoError(ths.T(), err)
			err = rfc2865.NASIPAddress_Set(packet, net.ParseIP(nas.Ip))
			require.NoError(ths.T(), err)
			err = rfc2865.NASIdentifier_AddString(packet, "local")
			require.NoError(ths.T(), err)
			err = rfc2865.NASPortType_Add(packet, rfc2865.NASPortType_Value_Ethernet)
			require.NoError(ths.T(), err)
			err = rfc2865.NASPort_Add(packet, rfc2865.NASPort(data.NasPort))
			require.NoError(ths.T(), err)

			err = rfc2866.AcctAuthentic_Add(packet, rfc2866.AcctAuthentic_Value_RADIUS)
			require.NoError(ths.T(), err)
			err = rfc2866.AcctSessionID_Add(packet, []byte(sessionId))
			require.NoError(ths.T(), err)
			err = rfc2866.AcctSessionTime_Set(packet, 0)
			require.NoError(ths.T(), err)

			radResponse, err = rad.Exchange(context.Background(), packet, "127.0.0.1:1813")
			require.NoError(ths.T(), err)
			require.NotNil(ths.T(), radResponse)

			ths.acctMap[c.Login] = data
			ip++
		}
	}

	ths.log.Info("Waiting for sessions created...")
	time.Sleep(SessionsCreateDelay)

	sessions, err := ths.sessionsRepo.Get()
	require.NoError(ths.T(), err)

	for _, a := range ths.acctMap {
		found := false
		for _, s := range sessions {
			if s.Login == a.Login {
				found = true
				break
			}
		}
		assert.Equal(ths.T(), true, found, "Session for %s not exists", a.Login)
	}

	require.Equal(ths.T(), len(customers)-len(shouldReject), len(sessions))
}

func (ths *RadiusServerTestSuite) TestRadius_03_Attrs() {
	err := ths.sessionsRepo.DeleteAll()
	require.NoError(ths.T(), err)
	err = ths.sessionCache.DeleteAll()
	require.NoError(ths.T(), err)

	if ths.radiusEventHandler == nil {
		ths.radiusEventHandler, err = radius.NewRadiusEventHandler(
			ths.cfg,
			ths.log,
			ths.customersRepo,
			ths.nasesRepo,
			ths.leasesRepo,
			ths.sessionsRepo,
			ths.sessionsLogRepo,
			ths.serviceInternetRepo,
		)
		require.NoError(ths.T(), err)

		go ths.radiusEventHandler.Start()

		ths.radiusEventHandler.Wg.Wait()
		ths.log.Info("RadiusEventHandler ready")
	}

	if ths.radiusServer == nil {
		rsEvents, err := events.NewEvents(ths.cfg, ths.log, uuid.NewString(), radius.Queue)
		require.NoError(ths.T(), err)

		radiusServer := radius.NewRadiusServer(ths.cfg, ths.log, rsEvents)

		go func() {
			err := radiusServer.Start()
			require.NoError(ths.T(), err)
		}()

		radiusServer.Wg.Wait()
		ths.log.Info("RadiusServer ready")
	}

	nases, err := ths.nasesRepo.Get()
	nas := nases[0]

	packet := rad.New(rad.CodeAccessRequest, []byte(nas.Secret))
	err = rfc2865.UserName_SetString(packet, fixtures.FakeCustomers()[0].Login)
	require.NoError(ths.T(), err)
	err = rfc2865.UserPassword_SetString(packet, fixtures.FakeCustomers()[0].Password)
	require.NoError(ths.T(), err)
	err = rfc2865.NASIPAddress_Set(packet, net.ParseIP(nas.Ip))
	require.NoError(ths.T(), err)
	err = rfc2865.NASIdentifier_SetString(packet, "local")
	require.NoError(ths.T(), err)
	err = rfc2865.NASPortType_Set(packet, rfc2865.NASPortType_Value_Ethernet)
	require.NoError(ths.T(), err)
	radResponse, err := rad.Exchange(context.Background(), packet, "127.0.0.1:1812")
	assert.NoError(ths.T(), err)
	assert.NotNil(ths.T(), radResponse)

	if radResponse.Code == rad.CodeAccessAccept {
		framedPool := rfc2869.FramedPool_Get(radResponse)
		assert.Equal(ths.T(), "og_pool", string(framedPool))
		idleTimeout := rfc2865.IdleTimeout_Get(radResponse)
		assert.Equal(ths.T(), rfc2865.IdleTimeout(600), idleTimeout)
	}
}

func TestRadiusServerTestSuite(t *testing.T) {
	testingSuite := new(RadiusServerTestSuite)
	testingSuite.eventsEnabled = true
	suite.Run(t, testingSuite)
}
