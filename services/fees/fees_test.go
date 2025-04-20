package fees

import (
	fixtures2 "code.evixo.ru/ncc/ncc-backend/pkg/fixtures"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFees(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := zero.NewLogger()

	ShouldBeProcessed := 16
	ShouldBeBlocked := 6
	ShouldSetDays := 2

	feeRepo := mocks.NewMockFees(ctrl)
	feeRepo.EXPECT().GetProcessedMap(gomock.Any()).Times(1).Return(fixtures2.FakeFeeProcessedMap(), nil)
	feeRepo.EXPECT().Create(gomock.Any()).Times(ShouldBeProcessed).Return(nil)

	processed, err := feeRepo.GetProcessedMap()
	assert.NoError(t, err)

	customerRepo := mocks.NewMockCustomers(ctrl)
	customerRepo.EXPECT().SetDeposit(gomock.Any(), gomock.Any()).Times(ShouldBeProcessed).Return(nil)
	customerRepo.EXPECT().SetServiceInternetState(gomock.Any(), gomock.Any()).Times(ShouldBeBlocked).Return(nil)
	customerRepo.EXPECT().SetCreditDaysLeft(gomock.Any(), gomock.Any()).Times(ShouldSetDays).Return(nil)

	internetMap := map[string]models2.ServiceInternetData{}
	for _, i := range fixtures2.FakeServicesInternet() {
		internetMap[i.Id] = i
	}

	customDataMap := map[string]models2.ServiceInternetCustomData{}
	for _, i := range fixtures2.FakeServiceInternetCustomData() {
		customDataMap[i.CustomerId.UUID.String()] = i
	}

	serviceInternetRepo := mocks.NewMockServiceInternet(ctrl)
	ipPoolRepo := mocks.NewMockIpPools(ctrl)

	ipPoolRepo.EXPECT().Get().AnyTimes().Return(fixtures2.FakeIpPools(), nil)

	fees := NewFees(log, feeRepo, customerRepo, serviceInternetRepo, ipPoolRepo)
	days := fees.DaysIn(9, 2022)
	assert.Equal(t, days, 30)

	forTime := time.Date(2022, 9, 8, 4, 0, 0, 0, time.Now().Location())
	feeDatas, err := fees.Process(
		internetMap,
		fixtures2.FakeCustomers(),
		customDataMap,
		processed,
		30,
		forTime,
		false,
		0,
		false,
		false)

	assert.NoError(t, err)
	assert.Len(t, feeDatas, ShouldBeProcessed)

	assert.Equal(t, 90.0, feeDatas["user1"].NewDeposit)
	assert.Equal(t, 10.00, feeDatas["user1"].FeeLog.FeeAmount)
	assert.Equal(t, models2.CustomerStateActive, feeDatas["user1"].NewState)

	assert.Equal(t, -5.0, feeDatas["user3"].NewDeposit)
	assert.Equal(t, 10.00, feeDatas["user3"].FeeLog.FeeAmount)
	assert.Equal(t, models2.CustomerStateActive, feeDatas["user3"].NewState)

	assert.Equal(t, -5.0, feeDatas["user4"].NewDeposit)
	assert.Equal(t, 0.00, feeDatas["user4"].FeeLog.FeeAmount)
	assert.Equal(t, models2.CustomerStateBlocked, feeDatas["user4"].NewState)

	assert.Equal(t, 88.0, feeDatas["user5"].NewDeposit)
	assert.Equal(t, 12.00, feeDatas["user5"].FeeLog.FeeAmount)
	assert.Equal(t, models2.CustomerStateActive, feeDatas["user5"].NewState)

	assert.Equal(t, 0.0003, feeDatas["user6"].NewDeposit)
	assert.Equal(t, 0.00, feeDatas["user6"].FeeLog.FeeAmount)
	assert.Equal(t, models2.CustomerStateBlocked, feeDatas["user6"].NewState)

	assert.Equal(t, -15.0, feeDatas["user7"].NewDeposit)
	assert.Equal(t, 10.00, feeDatas["user7"].FeeLog.FeeAmount)
	assert.Equal(t, 1, feeDatas["user7"].NewCreditDaysLeft)
	assert.Equal(t, models2.CustomerStateActive, feeDatas["user7"].NewState)

	assert.Equal(t, -15.0, feeDatas["user8"].NewDeposit)
	assert.Equal(t, 10.00, feeDatas["user8"].FeeLog.FeeAmount)
	assert.Equal(t, 0, feeDatas["user8"].NewCreditDaysLeft)
	assert.Equal(t, models2.CustomerStateActive, feeDatas["user8"].NewState)

	assert.Equal(t, -15.0, feeDatas["user9"].NewDeposit)
	assert.Equal(t, 10.00, feeDatas["user9"].FeeLog.FeeAmount)
	assert.Equal(t, models2.CustomerStateActive, feeDatas["user9"].NewState)

	assert.Equal(t, -5.0, feeDatas["user10"].NewDeposit)
	assert.Equal(t, 0.00, feeDatas["user10"].FeeLog.FeeAmount)
	assert.Equal(t, models2.CustomerStateBlocked, feeDatas["user10"].NewState)

	assert.Equal(t, -15.0, feeDatas["user11"].NewDeposit)
	assert.Equal(t, 10.00, feeDatas["user11"].FeeLog.FeeAmount)
	assert.Equal(t, models2.CustomerStateActive, feeDatas["user11"].NewState)

	assert.Equal(t, -15.0, feeDatas["user12"].NewDeposit)
	assert.Equal(t, 0.00, feeDatas["user12"].FeeLog.FeeAmount)
	assert.Equal(t, models2.CustomerStateBlocked, feeDatas["user12"].NewState)

	assert.Equal(t, 0.00, feeDatas["user13"].NewDeposit)
	assert.Equal(t, 0.00, feeDatas["user13"].FeeLog.FeeAmount)
	assert.Equal(t, models2.CustomerStateBlocked, feeDatas["user13"].NewState)

	assert.Equal(t, -9.99, feeDatas["user14"].NewDeposit)
	assert.Equal(t, 10.00, feeDatas["user14"].FeeLog.FeeAmount)
	assert.Equal(t, models2.CustomerStateActive, feeDatas["user14"].NewState)

	assert.Equal(t, 0.001, feeDatas["user15"].NewDeposit)
	assert.Equal(t, 0.00, feeDatas["user15"].FeeLog.FeeAmount)
	assert.Equal(t, models2.CustomerStateBlocked, feeDatas["user15"].NewState)

	assert.Empty(t, feeDatas["user16"])

	assert.Empty(t, feeDatas["user20"])
}
