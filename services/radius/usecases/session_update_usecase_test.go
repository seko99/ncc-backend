package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/fixtures"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository/mocks"
	"code.evixo.ru/ncc/ncc-backend/services/radius/dto"
	mocks2 "code.evixo.ru/ncc/ncc-backend/services/radius/mocks"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type SessionUpdateTestSuite struct {
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

func (ths *SessionUpdateTestSuite) SetupSuite() {
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

func (ths *SessionUpdateTestSuite) TestSessionUpdateUsecase_Execute() {

	fakeSession := fixtures.FakeSessions()[0]
	ths.sessions.EXPECT().GetBySessionId(gomock.Any()).Times(1).Return(fakeSession, nil)
	ths.sessions.EXPECT().UpdateBySessionId(gomock.Any()).Times(1).Return(nil)

	usecase := NewSessionUpdateUsecase(
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
		ths.startUsecase,
	)

	response, err := usecase.Execute(dto.SessionUpdateRequest{
		AcctSessionId: fakeSession.AcctSessionId,
		UserName:      fakeSession.Ip,
		FramedIp:      fakeSession.Ip,
		NasIpAddress:  fakeSession.Nas.Ip,
	})
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), response)
}

func (ths *SessionUpdateTestSuite) TestSessionUpdateUsecase_SessionNotExists_NotInLog_Execute() {

	fakeSession := models.SessionData{}
	ths.sessions.EXPECT().GetBySessionId(gomock.Any()).Times(1).Return(fakeSession, gorm.ErrRecordNotFound)

	fakeSessionLog := models.SessionsLogData{}
	ths.sessionsLog.EXPECT().GetBySessionId(gomock.Any()).Times(1).Return(fakeSessionLog, gorm.ErrRecordNotFound)

	fakeStartResponse := dto.SessionStartResponse{
		AcctSessionId: uuid.NewString(),
	}
	ths.startUsecase.EXPECT().Execute(gomock.Any()).Times(1).Return(&fakeStartResponse, nil)

	usecase := NewSessionUpdateUsecase(
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
		ths.startUsecase,
	)

	response, err := usecase.Execute(dto.SessionUpdateRequest{
		AcctSessionId: uuid.NewString(),
		UserName:      "user1",
		FramedIp:      "10.50.0.11",
		NasIpAddress:  "10.0.27.1",
	})
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), response)
}

func (ths *SessionUpdateTestSuite) TestSessionUpdateUsecase_SessionNotExists_ExistsInLog_Execute() {

	fakeSession := models.SessionData{}
	ths.sessions.EXPECT().GetBySessionId(gomock.Any()).Times(1).Return(fakeSession, gorm.ErrRecordNotFound)
	ths.sessions.EXPECT().Create(gomock.Any()).Times(1).Return(nil)

	sessionId := uuid.NewString()
	fakeSessionLog := models.SessionsLogData{
		AcctSessionId: sessionId,
	}
	ths.sessionsLog.EXPECT().GetBySessionId(gomock.Any()).Times(1).Return(fakeSessionLog, nil)
	ths.sessionsLog.EXPECT().DeleteById(gomock.Any()).Times(1).Return(nil)

	usecase := NewSessionUpdateUsecase(
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
		ths.startUsecase,
	)

	response, err := usecase.Execute(dto.SessionUpdateRequest{
		AcctSessionId: uuid.NewString(),
		UserName:      "user1",
		FramedIp:      "10.50.0.11",
		NasIpAddress:  "10.0.27.1",
	})
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), response)
}

func TestSessionUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(SessionUpdateTestSuite))
}
