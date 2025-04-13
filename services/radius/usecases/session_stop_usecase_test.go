package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/fixtures"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository/mocks"
	"code.evixo.ru/ncc/ncc-backend/services/radius/dto"
	mocks2 "code.evixo.ru/ncc/ncc-backend/services/radius/mocks"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SessionStopTestSuite struct {
	suite.Suite
	ctx  context.Context
	log  logger.Logger
	cfg  *config.Config
	ctrl *gomock.Controller

	customers    *mocks.MockCustomers
	sessions     *mocks.MockSessions
	sessionsLog  *mocks.MockSessionsLog
	leases       *mocks.MockDhcpLeases
	nases        *mocks.MockNases
	nasTypes     *mocks.MockNasTypes
	startUsecase *mocks2.MockSessionStartUsecase
}

func (ths *SessionStopTestSuite) SetupSuite() {
	ths.ctrl = gomock.NewController(ths.T())

	ths.ctx = context.Background()
	ths.log = zero.NewLogger()

	ths.sessions = mocks.NewMockSessions(ths.ctrl)
	ths.sessionsLog = mocks.NewMockSessionsLog(ths.ctrl)
	ths.leases = mocks.NewMockDhcpLeases(ths.ctrl)
	ths.customers = mocks.NewMockCustomers(ths.ctrl)
	ths.nases = mocks.NewMockNases(ths.ctrl)

	ths.startUsecase = mocks2.NewMockSessionStartUsecase(ths.ctrl)
}

func (ths *SessionStopTestSuite) TestSessionStopUsecase_Execute() {

	fakeSession := fixtures.FakeSessions()[0]
	ths.sessions.EXPECT().GetBySessionId(gomock.Any()).Times(1).Return(fakeSession, nil)
	ths.sessions.EXPECT().UpdateBySessionId(gomock.Any()).Times(1).Return(nil)
	ths.sessions.EXPECT().DeleteBySessionId(gomock.Any()).Times(1).Return(nil)

	ths.sessionsLog.EXPECT().Create(gomock.Any()).Times(1).Return(nil)

	usecase := NewSessionStopUsecase(
		ths.cfg,
		ths.log,
		ths.sessions,
		ths.sessionsLog,
		ths.leases,
		ths.customers,
		ths.nases,
		ths.sessions,
		ths.leases,
		ths.customers,
		ths.nases,
	)

	response, err := usecase.Execute(dto.SessionStopRequest{
		AcctSessionId: fakeSession.AcctSessionId,
		UserName:      fakeSession.Ip,
		FramedIp:      fakeSession.Ip,
		NasIpAddress:  fakeSession.Nas.Ip,
	})
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), response)
}

func TestSessionStopTestSuite(t *testing.T) {
	suite.Run(t, new(SessionStopTestSuite))
}
