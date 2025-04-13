package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/fixtures"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository/memory"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository/mocks"
	"code.evixo.ru/ncc/ncc-backend/services/radius/dto"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SessionStartTestSuite struct {
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
}

func (ths *SessionStartTestSuite) SetupSuite() {
	ths.ctrl = gomock.NewController(ths.T())

	ths.ctx = context.Background()
	ths.log = zero.NewLogger()

	ths.sessions = mocks.NewMockSessions(ths.ctrl)
	ths.sessionsLog = mocks.NewMockSessionsLog(ths.ctrl)
	ths.leases = mocks.NewMockDhcpLeases(ths.ctrl)
	ths.customers = mocks.NewMockCustomers(ths.ctrl)
	ths.nases = mocks.NewMockNases(ths.ctrl)
}

func (ths *SessionStartTestSuite) TestSessionStartUsecase_Correct_By_Login_Execute() {
	ths.sessions.EXPECT().Create(gomock.Any()).Times(1).Return(nil)
	ths.sessions.EXPECT().GetByLogin(gomock.Any()).Times(1).Return(nil, memory.ErrNotFound)

	fakeCustomer := fixtures.FakeCustomers()[3]
	ths.customers.EXPECT().GetByLogin(gomock.Any()).Times(1).Return(&fakeCustomer, nil)

	fakeNas := fixtures.FakeNases()[0]
	ths.nases.EXPECT().GetByIP(gomock.Any()).Times(1).Return(fakeNas, nil)

	usecase := NewSessionStartUsecase(ths.cfg, ths.log, ths.sessions, ths.sessionsLog, ths.leases, ths.customers, ths.nases, ths.sessions, ths.leases, ths.customers, ths.nases)

	response, err := usecase.Execute(dto.SessionStartRequest{
		AcctSessionId: uuid.NewString(),
		UserName:      "user4",
		FramedIp:      "10.50.0.13",
		NasIpAddress:  "10.0.27.1",
	})
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), response)
}

func (ths *SessionStartTestSuite) TestSessionStartUsecase_Incorrect_By_Login_Duplicate_Execute() {
	fakeSessions := fixtures.FakeSessions()
	ths.sessions.EXPECT().GetByLogin(gomock.Any()).Times(1).Return(fakeSessions, nil)

	usecase := NewSessionStartUsecase(ths.cfg, ths.log, ths.sessions, ths.sessionsLog, ths.leases, ths.customers, ths.nases, ths.sessions, ths.leases, ths.customers, ths.nases)

	response, err := usecase.Execute(dto.SessionStartRequest{
		AcctSessionId: fixtures.FakeSessions()[0].AcctSessionId,
		UserName:      "user1",
		FramedIp:      "10.50.0.11",
		NasIpAddress:  "10.0.27.1",
	})
	assert.ErrorIs(ths.T(), err, ErrDuplicate)
	assert.Nil(ths.T(), response)
}

func (ths *SessionStartTestSuite) TestSessionStartUsecase_Incorrect_By_Login_MaxSessions_Execute() {
	fakeSessions := fixtures.FakeSessions()
	ths.sessions.EXPECT().GetByLogin(gomock.Any()).Times(1).Return(fakeSessions, nil)

	usecase := NewSessionStartUsecase(ths.cfg, ths.log, ths.sessions, ths.sessionsLog, ths.leases, ths.customers, ths.nases, ths.sessions, ths.leases, ths.customers, ths.nases)

	response, err := usecase.Execute(dto.SessionStartRequest{
		AcctSessionId: uuid.NewString(),
		UserName:      "user1",
		FramedIp:      "10.50.0.13",
		NasIpAddress:  "10.0.27.1",
	})
	assert.ErrorIs(ths.T(), err, ErrMaxSessions)
	assert.Nil(ths.T(), response)
}

func (ths *SessionStartTestSuite) TestSessionStartUsecase_Correct_By_IP_Execute() {
	ths.sessions.EXPECT().Create(gomock.Any()).Times(1).Return(nil)
	ths.sessions.EXPECT().GetByIP(gomock.Any()).Times(1).Return(models.SessionData{}, ErrNotFound)

	fakeLease := fixtures.FakeLeases()[2]
	ths.leases.EXPECT().GetByIP(gomock.Any()).Times(1).Return(fakeLease, nil)

	fakeNas := fixtures.FakeNases()[0]
	ths.nases.EXPECT().GetByIP(gomock.Any()).Times(1).Return(fakeNas, nil)

	usecase := NewSessionStartUsecase(ths.cfg, ths.log, ths.sessions, ths.sessionsLog, ths.leases, ths.customers, ths.nases, ths.sessions, ths.leases, ths.customers, ths.nases)

	response, err := usecase.Execute(dto.SessionStartRequest{
		AcctSessionId: uuid.NewString(),
		UserName:      "10.50.0.13",
		FramedIp:      "10.50.0.13",
		NasIpAddress:  "10.0.27.1",
	})
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), response)
}

func (ths *SessionStartTestSuite) TestSessionStartUsecase_Inorrect_By_IP_NoLease_Execute() {
	ths.sessions.EXPECT().GetByIP(gomock.Any()).Times(1).Return(models.SessionData{}, ErrNotFound)

	fakeLease := models.LeaseData{}
	ths.leases.EXPECT().GetByIP(gomock.Any()).Times(1).Return(fakeLease, ErrNotFound)

	usecase := NewSessionStartUsecase(ths.cfg, ths.log, ths.sessions, ths.sessionsLog, ths.leases, ths.customers, ths.nases, ths.sessions, ths.leases, ths.customers, ths.nases)

	response, err := usecase.Execute(dto.SessionStartRequest{
		AcctSessionId: uuid.NewString(),
		UserName:      "10.50.0.13",
		FramedIp:      "10.50.0.13",
		NasIpAddress:  "10.0.27.1",
	})
	assert.ErrorIs(ths.T(), err, ErrNoLease)
	assert.Nil(ths.T(), response)
}

func (ths *SessionStartTestSuite) TestSessionStartUsecase_Inorrect_By_IP_NoBinding_Execute() {
	ths.sessions.EXPECT().GetByIP(gomock.Any()).Times(1).Return(models.SessionData{}, ErrNotFound)

	fakeLease := fixtures.FakeLeases()[3]
	ths.leases.EXPECT().GetByIP(gomock.Any()).Times(1).Return(fakeLease, nil)

	usecase := NewSessionStartUsecase(ths.cfg, ths.log, ths.sessions, ths.sessionsLog, ths.leases, ths.customers, ths.nases, ths.sessions, ths.leases, ths.customers, ths.nases)

	response, err := usecase.Execute(dto.SessionStartRequest{
		AcctSessionId: uuid.NewString(),
		UserName:      "10.50.0.14",
		FramedIp:      "10.50.0.14",
		NasIpAddress:  "10.0.27.1",
	})
	assert.ErrorIs(ths.T(), err, ErrNoBinding)
	assert.Nil(ths.T(), response)
}

func (ths *SessionStartTestSuite) TestSessionStartUsecase_Inorrect_By_IP_NotEqual_Execute() {
	usecase := NewSessionStartUsecase(ths.cfg, ths.log, ths.sessions, ths.sessionsLog, ths.leases, ths.customers, ths.nases, ths.sessions, ths.leases, ths.customers, ths.nases)

	response, err := usecase.Execute(dto.SessionStartRequest{
		AcctSessionId: uuid.NewString(),
		UserName:      "10.50.0.13",
		FramedIp:      "10.50.0.18",
		NasIpAddress:  "10.0.27.1",
	})
	assert.ErrorIs(ths.T(), err, ErrIPNotEqual)
	assert.Nil(ths.T(), response)
}

func (ths *SessionStartTestSuite) TestSessionStartUsecase_Incorrect_By_IP_Duplicate_Execute() {
	fakeSession := fixtures.FakeSessions()[0]
	ths.sessions.EXPECT().GetByIP(gomock.Any()).Times(1).Return(fakeSession, nil)

	usecase := NewSessionStartUsecase(ths.cfg, ths.log, ths.sessions, ths.sessionsLog, ths.leases, ths.customers, ths.nases, ths.sessions, ths.leases, ths.customers, ths.nases)

	response, err := usecase.Execute(dto.SessionStartRequest{
		AcctSessionId: fakeSession.AcctSessionId,
		UserName:      "10.50.0.11",
		FramedIp:      "10.50.0.11",
		NasIpAddress:  "10.0.27.1",
	})
	assert.ErrorIs(ths.T(), err, ErrDuplicate)
	assert.Nil(ths.T(), response)
}

func (ths *SessionStartTestSuite) TestSessionStartUsecase_Incorrect_By_IP_NewSession_Execute() {
	fakeSession := fixtures.FakeSessions()[0]
	ths.sessions.EXPECT().GetByIP(gomock.Any()).Times(1).Return(fakeSession, nil)
	ths.sessions.EXPECT().Delete(gomock.Any()).Times(1).Return(nil)
	ths.sessions.EXPECT().Create(gomock.Any()).Times(1).Return(nil)

	ths.sessionsLog.EXPECT().Create(gomock.Any()).Times(1).Return(nil)

	fakeNas := fixtures.FakeNases()[0]
	ths.nases.EXPECT().GetByIP(gomock.Any()).Times(1).Return(fakeNas, nil)

	usecase := NewSessionStartUsecase(ths.cfg, ths.log, ths.sessions, ths.sessionsLog, ths.leases, ths.customers, ths.nases, ths.sessions, ths.leases, ths.customers, ths.nases)

	response, err := usecase.Execute(dto.SessionStartRequest{
		AcctSessionId: uuid.NewString(),
		UserName:      "10.50.0.11",
		FramedIp:      "10.50.0.11",
		NasIpAddress:  "10.0.27.1",
	})
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), response)
}

func TestSessionStartTestSuite(t *testing.T) {
	suite.Run(t, new(SessionStartTestSuite))
}
