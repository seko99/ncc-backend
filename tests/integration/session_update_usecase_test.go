package integration

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/fixtures"
	"code.evixo.ru/ncc/ncc-backend/services/radius/dto"
	"code.evixo.ru/ncc/ncc-backend/services/radius/usecases"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SessionUpdateTestSuite struct {
	BaseTestSuite
}

func (ths *SessionUpdateTestSuite) TestSessionUpdateUsecase_SessionExists_Execute() {
	ucSessionStart := usecases.NewSessionStartUsecase(
		ths.cfg,
		ths.log,
		ths.sessionsRepo,
		ths.sessionsLogRepo,
		ths.leasesRepo,
		ths.customersRepo,
		ths.nasesRepo,
		ths.sessionCache,
		ths.leasesCache,
		ths.customersCache,
		ths.nasesCache,
	)

	usecase := usecases.NewSessionUpdateUsecase(
		ths.cfg,
		ths.log,
		ths.sessionsRepo,
		ths.sessionsLogRepo,
		ths.leasesRepo,
		ths.customersRepo,
		ths.nasesRepo,
		ths.sessionCache,
		ths.leasesCache,
		ths.customersCache,
		ths.nasesCache,
		&ucSessionStart,
	)

	fakeSession := fixtures.FakeSessions()[0]

	response, err := usecase.Execute(dto.SessionUpdateRequest{
		AcctSessionId:    fakeSession.AcctSessionId,
		UserName:         fakeSession.Ip,
		FramedIp:         fakeSession.Ip,
		NasIpAddress:     "10.0.27.1",
		AcctSessionTime:  fakeSession.Duration,
		AcctInputOctets:  uint32(fakeSession.OctetsIn),
		AcctOutputOctets: uint32(fakeSession.OctetsOut),
	})
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), response)

	updatedSession, err := ths.sessionsRepo.GetByLogin(fakeSession.Login)
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), updatedSession)
	assert.Len(ths.T(), updatedSession, 1)
	assert.Equal(ths.T(), updatedSession[0].AcctSessionId, response.AcctSessionId)
	assert.Equal(ths.T(), updatedSession[0].Ip, fakeSession.Ip)
	assert.Equal(ths.T(), updatedSession[0].Duration, fakeSession.Duration)
	assert.Equal(ths.T(), updatedSession[0].OctetsIn, fakeSession.OctetsIn)
	assert.Equal(ths.T(), updatedSession[0].OctetsOut, fakeSession.OctetsOut)
}

func (ths *SessionUpdateTestSuite) TestSessionUpdateUsecase_SessionNotExists_ExistsInLog_Execute() {
	ucSessionStart := usecases.NewSessionStartUsecase(
		ths.cfg,
		ths.log,
		ths.sessionsRepo,
		ths.sessionsLogRepo,
		ths.leasesRepo,
		ths.customersRepo,
		ths.nasesRepo,
		ths.sessionCache,
		ths.leasesCache,
		ths.customersCache,
		ths.nasesCache,
	)

	usecase := usecases.NewSessionUpdateUsecase(
		ths.cfg,
		ths.log,
		ths.sessionsRepo,
		ths.sessionsLogRepo,
		ths.leasesRepo,
		ths.customersRepo,
		ths.nasesRepo,
		ths.sessionCache,
		ths.leasesCache,
		ths.customersCache,
		ths.nasesCache,
		&ucSessionStart,
	)

	user := "user4"
	ip := "10.50.0.15"

	response, err := usecase.Execute(dto.SessionUpdateRequest{
		AcctSessionId: uuid.NewString(),
		UserName:      user,
		FramedIp:      ip,
		NasIpAddress:  "127.0.0.1",
	})
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), response)

	createdSession, err := ths.sessionsRepo.GetByLogin(user)
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), createdSession)
	assert.Len(ths.T(), createdSession, 1)
	assert.Equal(ths.T(), createdSession[0].AcctSessionId, response.AcctSessionId)
	assert.Equal(ths.T(), createdSession[0].Ip, ip)
}

func TestSessionUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(SessionUpdateTestSuite))
}
