package integration

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/fixtures"
	"code.evixo.ru/ncc/ncc-backend/services/dhcp"
	"github.com/gogf/gf/net/gipv4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type DhcpServerTestSuite struct {
	BaseTestSuite
	dhcpEventHandler *dhcp.EventHandler
	dhcpServer       *dhcp.Server
	rsEvents         *events.Events
}

func (ths *DhcpServerTestSuite) TestDhcp_00_InitConfig() {
	var err error

	ths.dhcpEventHandler, err = dhcp.NewDhcpEventHandler(
		ths.cfg,
		ths.log,
		ths.nasesRepo,
		ths.leasesRepo,
		ths.poolsRepo,
		ths.bindingsRepo,
	)
	require.NoError(ths.T(), err)

	go ths.dhcpEventHandler.Start()

	ths.dhcpEventHandler.Wg.Wait()
	ths.log.Info("DhcpEventHandler ready")

	ths.rsEvents, err = events.NewEvents(ths.cfg, ths.log, uuid.NewString(), dhcp.Queue)
	require.NoError(ths.T(), err)

	ths.dhcpServer = dhcp.NewDhcpServer(ths.cfg, ths.log, ths.rsEvents)

	go func() {
		err := ths.dhcpServer.Start()
		require.NoError(ths.T(), err)
	}()

	ths.dhcpServer.Wg.Wait()
	ths.log.Info("DhcpServer ready")

	client, err := dhcp.NewDhcpClient("127.0.0.1:1067")
	assert.NoError(ths.T(), err)
	assert.NotNil(ths.T(), client)

	poolStart := gipv4.Ip2long(fixtures.FakePools()[1].RangeStart)
	poolEnd := gipv4.Ip2long(fixtures.FakePools()[1].RangeEnd)
	pkt := dhcp.Packet{}

	// Discover 01:02:03:04:05:06
	buf, err := pkt.Discover("01:02:03:04:05:06")
	assert.NoError(ths.T(), err)
	reply, err := client.SendPacket(buf)
	assert.NotNil(ths.T(), reply)
	assert.Equal(ths.T(), byte(dhcp.MsgTypeOffer), reply.Options.Opt53.Type)
	assert.GreaterOrEqual(ths.T(), reply.Packet.Ciaddr, poolStart)
	assert.LessOrEqual(ths.T(), reply.Packet.Ciaddr, poolEnd)

	offeredIp := reply.Packet.Ciaddr

	// Discover 01:02:03:04:05:06 again - get same IP
	buf, err = pkt.Discover("01:02:03:04:05:06")
	assert.NoError(ths.T(), err)
	reply, err = client.SendPacket(buf)
	assert.NotNil(ths.T(), reply)
	assert.Equal(ths.T(), byte(dhcp.MsgTypeOffer), reply.Options.Opt53.Type)
	assert.Equal(ths.T(), offeredIp, reply.Packet.Ciaddr)

	// Discover 06:05:04:03:02:01 - get next IP
	buf, err = pkt.Discover("06:05:04:03:02:01")
	assert.NoError(ths.T(), err)
	reply, err = client.SendPacket(buf)
	assert.NotNil(ths.T(), reply)
	assert.Equal(ths.T(), byte(dhcp.MsgTypeOffer), reply.Options.Opt53.Type)
	assert.Equal(ths.T(), offeredIp+1, reply.Packet.Ciaddr)

	// Request 01:02:03:04:05:06 - ACK, status=Accepted
	buf, err = pkt.Request("01:02:03:04:05:06", offeredIp)
	assert.NoError(ths.T(), err)
	reply, err = client.SendPacket(buf)
	assert.NotNil(ths.T(), reply)
	assert.Equal(ths.T(), byte(dhcp.MsgTypeAck), reply.Options.Opt53.Type)
	assert.Equal(ths.T(), offeredIp, reply.Packet.Ciaddr)
	assert.Equal(ths.T(), true, ths.dhcpServer.IsAllocated(offeredIp))
	assert.Contains(ths.T(), ths.dhcpServer.GetServerLeases(), offeredIp)
	lease, err := ths.dhcpServer.GetServerLeaseByIP(offeredIp)
	assert.NoError(ths.T(), err)
	assert.Equal(ths.T(), dhcp.LeaseStatusAccepted, int(lease.Status))
	prevExpire := lease.Expire

	// wait at least 1 second or newExpire will be same as prevExpire (now()+pool.LeaseTime)
	time.Sleep(1 * time.Second)

	// Request 01:02:03:04:05:06 - ACK, expire > prevExpire (renew)
	buf, err = pkt.Request("01:02:03:04:05:06", offeredIp)
	assert.NoError(ths.T(), err)
	reply, err = client.SendPacket(buf)
	require.NotNil(ths.T(), reply)
	require.Equal(ths.T(), byte(dhcp.MsgTypeAck), reply.Options.Opt53.Type)
	require.Equal(ths.T(), offeredIp, reply.Packet.Ciaddr)
	lease, err = ths.dhcpServer.GetServerLeaseByIP(offeredIp)
	assert.NoError(ths.T(), err)
	assert.Greater(ths.T(), lease.Expire, prevExpire)

	// Request 00:00:00:00:00:00 - NAK, not exists
	buf, err = pkt.Request("00:00:00:00:00:00", 0)
	assert.NoError(ths.T(), err)
	reply, err = client.SendPacket(buf)
	assert.NotNil(ths.T(), reply)
	assert.Equal(ths.T(), byte(dhcp.MsgTypeNak), reply.Options.Opt53.Type)

	err = client.Close()
	assert.NoError(ths.T(), err)
}

func TestDhcpServerTestSuite(t *testing.T) {
	testingSuite := new(DhcpServerTestSuite)
	testingSuite.eventsEnabled = true
	suite.Run(t, testingSuite)
}
