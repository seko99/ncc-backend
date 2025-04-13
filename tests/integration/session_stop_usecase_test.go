package integration

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/fixtures"
	"code.evixo.ru/ncc/ncc-backend/services/radius/dto"
	"code.evixo.ru/ncc/ncc-backend/services/radius/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SessionStopTestSuite struct {
	BaseTestSuite
}

func (ths *SessionStopTestSuite) TestSessionStopUsecase_Execute() {
	usecase := usecases.NewSessionStopUsecase(
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

	fakeSession := fixtures.FakeSessions()[0]

	response, err := usecase.Execute(dto.SessionStopRequest{
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

	stoppedSession, err := ths.sessionsLogRepo.GetBySessionId(fakeSession.AcctSessionId)
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), stoppedSession)
	assert.True(ths.T(), stoppedSession.StopTime.After(stoppedSession.StartTime))
	assert.Equal(ths.T(), fakeSession.Duration, stoppedSession.Duration)
	assert.Equal(ths.T(), fakeSession.OctetsIn, stoppedSession.OctetsIn)
	assert.Equal(ths.T(), fakeSession.OctetsOut, stoppedSession.OctetsOut)

	loggedSession, err := ths.sessionsLogRepo.GetBySessionId(fakeSession.AcctSessionId)
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), loggedSession)
}

func TestSessionStopTestSuite(t *testing.T) {
	suite.Run(t, new(SessionStopTestSuite))
}
