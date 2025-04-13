package integration

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/fixtures"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository/memory"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type DistributedCacheTestSuite struct {
	BaseTestSuite
}

const (
	Instances  = 5
	Operations = 10
)

func (ths *DistributedCacheTestSuite) TestSync() {
	var caches []repository.Sessions

	for i := 0; i < Instances; i++ {
		sessionCache, err := memory.NewSessions(ths.log, ths.sessionsRepo, ths.events)
		require.NoError(ths.T(), err)
		caches = append(caches, sessionCache)
	}

	for i := 0; i < Operations; i++ {
		id := uuid.NewString()
		sessionID := uuid.NewString()
		login := uuid.NewString()
		ip := faker.IPv4()
		newSession := models.SessionData{
			CommonData: models.CommonData{
				Id: id,
			},
			AcctSessionId:     sessionID,
			Login:             login,
			Ip:                ip,
			CustomerId:        models.NewNullUUID(fixtures.FakeCustomers()[0].Id),
			NasId:             models.NewNullUUID(fixtures.FakeNases()[0].Id),
			ServiceInternetId: models.NewNullUUID(fixtures.FakeServicesInternet()[0].Id),
		}
		err := ths.sessionsRepo.Create([]models.SessionData{newSession})
		require.NoError(ths.T(), err)

		time.Sleep(500 * time.Millisecond)

		for _, c := range caches {
			session, err := c.GetById(id)
			assert.NoError(ths.T(), err)
			assert.Equal(ths.T(), login, session.Login)
			assert.Equal(ths.T(), ip, session.Ip)
			assert.Equal(ths.T(), sessionID, session.AcctSessionId)
		}
	}
}

func TestDistributedCacheTestSuite(t *testing.T) {
	testingSuite := new(DistributedCacheTestSuite)
	testingSuite.eventsEnabled = true
	suite.Run(t, testingSuite)
}
