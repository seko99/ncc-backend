package integration

import (
	"code.evixo.ru/ncc/ncc-backend/services/radius/dto"
	"code.evixo.ru/ncc/ncc-backend/services/radius/usecases"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SessionStartTestSuite struct {
	BaseTestSuite
}

func (ths *SessionStartTestSuite) TestSessionStartUsecase_Correct_By_Login_Execute() {
	usecase := usecases.NewSessionStartUsecase(ths.cfg, ths.log, ths.sessionsRepo, ths.sessionsLogRepo, ths.leasesRepo, ths.customersRepo, ths.nasesRepo, ths.sessionCache, ths.leasesCache, ths.customersCache, ths.nasesCache)

	user := "user4"
	ip := "10.50.0.15"

	response, err := usecase.Execute(dto.SessionStartRequest{
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

func (ths *SessionStartTestSuite) TestSessionStartUsecase_Correct_By_IP_Execute() {
	usecase := usecases.NewSessionStartUsecase(ths.cfg, ths.log, ths.sessionsRepo, ths.sessionsLogRepo, ths.leasesRepo, ths.customersRepo, ths.nasesRepo, ths.sessionCache, ths.leasesCache, ths.customersCache, ths.nasesCache)

	ip := "10.50.0.13"

	response, err := usecase.Execute(dto.SessionStartRequest{
		AcctSessionId: uuid.NewString(),
		UserName:      ip,
		FramedIp:      ip,
		NasIpAddress:  "127.0.0.1",
	})
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), response)

	createdSession, err := ths.sessionsRepo.GetByIP(ip)
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), createdSession)
	assert.Equal(ths.T(), createdSession.AcctSessionId, response.AcctSessionId)
	assert.Equal(ths.T(), createdSession.Ip, ip)
}

func (ths *SessionStartTestSuite) TestSessionStartUsecase_Incorrect_By_IP_NewSession_Execute() {
	usecase := usecases.NewSessionStartUsecase(ths.cfg, ths.log, ths.sessionsRepo, ths.sessionsLogRepo, ths.leasesRepo, ths.customersRepo, ths.nasesRepo, ths.sessionCache, ths.leasesCache, ths.customersCache, ths.nasesCache)

	ip := "10.50.0.11"

	response, err := usecase.Execute(dto.SessionStartRequest{
		AcctSessionId: uuid.NewString(),
		UserName:      ip,
		FramedIp:      ip,
		NasIpAddress:  "127.0.0.1",
	})
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), response)

	createdSession, err := ths.sessionsRepo.GetByIP(ip)
	assert.NoError(ths.T(), err)
	assert.NotEmpty(ths.T(), createdSession)
	assert.Equal(ths.T(), createdSession.AcctSessionId, response.AcctSessionId)
	assert.Equal(ths.T(), createdSession.Ip, ip)
}

func (ths *SessionStartTestSuite) TestSessionStartUsecase_Incorrect_By_IP_NoBinding_Execute() {
	usecase := usecases.NewSessionStartUsecase(ths.cfg, ths.log, ths.sessionsRepo, ths.sessionsLogRepo, ths.leasesRepo, ths.customersRepo, ths.nasesRepo, ths.sessionCache, ths.leasesCache, ths.customersCache, ths.nasesCache)

	ip := "10.50.0.14"

	response, err := usecase.Execute(dto.SessionStartRequest{
		AcctSessionId: uuid.NewString(),
		UserName:      ip,
		FramedIp:      ip,
		NasIpAddress:  "127.0.0.1",
	})
	assert.ErrorIs(ths.T(), err, usecases.ErrNoBinding)
	assert.Empty(ths.T(), response)
}

func TestSessionStartTestSuite(t *testing.T) {
	suite.Run(t, new(SessionStartTestSuite))
}
