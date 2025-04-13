package integration

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/domain"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/providers/mocks"
	"code.evixo.ru/ncc/ncc-backend/services/informings"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type InformingsTestSuite struct {
	BaseTestSuite
}

func (ths *InformingsTestSuite) Test_01_List() {
	ctrl := gomock.NewController(ths.T())

	smsProvider := mocks.NewMockSmsProvider(ctrl)
	smsProvider.EXPECT().SendOne(gomock.Any(), "70000000016", gomock.Any()).Times(1).Return(nil)
	smsProvider.EXPECT().SendOne(gomock.Any(), "70000000017", gomock.Any()).Times(1).Return(nil)

	informingsService := informings.NewInformings(
		ths.log,
		smsProvider,
		ths.informingsRepo,
		ths.informingLogRepo,
		ths.informingsTestCustomersRepo,
		ths.customersRepo,
	)

	messageList, err := informingsService.PrepareMessageList()
	assert.NoError(ths.T(), err)
	assert.Equal(ths.T(), 2, len(messageList))
	assert.Equal(ths.T(), messageList[0].Phone, "70000000016")
	assert.Equal(ths.T(), messageList[1].Phone, "70000000017")

	err = informingsService.SendMessages(messageList)
	assert.NoError(ths.T(), err)

	log, err := ths.informingLogRepo.Get()
	assert.NoError(ths.T(), err)
	assert.Equal(ths.T(), 2, len(log))
	assert.Equal(ths.T(), "70000000016", log[0].Phone)
	assert.Equal(ths.T(), domain.MessageStatusSent, log[0].Status)
	assert.Equal(ths.T(), "70000000017", log[1].Phone)
	assert.Equal(ths.T(), domain.MessageStatusSent, log[1].Status)
}

func (ths *InformingsTestSuite) Test_02_ListAfterSent() {
	ctrl := gomock.NewController(ths.T())

	smsProvider := mocks.NewMockSmsProvider(ctrl)

	informingsService := informings.NewInformings(
		ths.log,
		smsProvider,
		ths.informingsRepo,
		ths.informingLogRepo,
		ths.informingsTestCustomersRepo,
		ths.customersRepo,
	)

	messageList, err := informingsService.PrepareMessageList()
	assert.NoError(ths.T(), err)
	assert.Equal(ths.T(), 0, len(messageList))
}

func (ths *InformingsTestSuite) Test_03_AfterSetFlag() {
	ctrl := gomock.NewController(ths.T())

	smsProvider := mocks.NewMockSmsProvider(ctrl)
	smsProvider.EXPECT().SendOne(gomock.Any(), "70000000017", gomock.Any()).Times(1).Return(nil)

	informingsService := informings.NewInformings(
		ths.log,
		smsProvider,
		ths.informingsRepo,
		ths.informingLogRepo,
		ths.informingsTestCustomersRepo,
		ths.customersRepo,
	)

	informingList, err := ths.informingsRepo.GetEnabled()
	assert.NoError(ths.T(), err)

	for _, i := range informingList {
		err := ths.informingsRepo.SetStart(i, time.Now().Add(-24*time.Hour))
		assert.NoError(ths.T(), err)
	}

	err = ths.customersRepo.SetFlag(models.CustomerData{}, models.CustomerFlagData{
		CustomerID: models.NewNullUUID("b3c3542d-c329-434d-91bd-e9fb83ba2878"),
		Name:       models.FieldSent,
		Val:        models.ExprFalse,
	})
	assert.NoError(ths.T(), err)

	messageList, err := informingsService.PrepareMessageList()
	assert.NoError(ths.T(), err)
	assert.Equal(ths.T(), 1, len(messageList))

	err = informingsService.SendMessages(messageList)
	assert.NoError(ths.T(), err)
}

func TestInformingsTestSuite(t *testing.T) {
	suite.Run(t, new(InformingsTestSuite))
}
