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

type SessionWatcherTestSuite struct {
	suite.Suite
	ctx  context.Context
	log  logger.Logger
	cfg  *config.Config
	ctrl *gomock.Controller

	customers   *mocks.MockCustomers
	sessions    *mocks.MockSessions
	sessionsLog *mocks.MockSessionsLog
	leases      *mocks.MockDhcpLeases
	nases       *mocks.MockNases
	nasTypes    *mocks.MockNasTypes
	stopUsecase *mocks2.MockSessionStopUsecase
}

func (ths *SessionWatcherTestSuite) SetupSuite() {
	ths.ctrl = gomock.NewController(ths.T())

	ths.ctx = context.Background()
	ths.log = zero.NewLogger()

	ths.sessions = mocks.NewMockSessions(ths.ctrl)
	ths.sessionsLog = mocks.NewMockSessionsLog(ths.ctrl)
	ths.leases = mocks.NewMockDhcpLeases(ths.ctrl)
	ths.customers = mocks.NewMockCustomers(ths.ctrl)
	ths.nases = mocks.NewMockNases(ths.ctrl)

	ths.stopUsecase = mocks2.NewMockSessionStopUsecase(ths.ctrl)
}

func (ths *SessionWatcherTestSuite) TestSessionWatcherUsecase_Execute() {

	fakeSessions := fixtures.FakeSessions()
	ths.sessions.EXPECT().Get().Times(1).Return(fakeSessions, nil)
	ths.stopUsecase.EXPECT().Execute(gomock.Any()).Times(1).Return(&dto.SessionStopResponse{}, nil)

	usecase := NewSessionWatcherUsecase(
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
		ths.stopUsecase,
	)

	err := usecase.Execute()
	assert.NoError(ths.T(), err)
}

func TestSessionWatcherTestSuite(t *testing.T) {
	suite.Run(t, new(SessionWatcherTestSuite))
}
