package informings

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/fixtures"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	"code.evixo.ru/ncc/ncc-backend/pkg/providers/mocks"
	mocks2 "code.evixo.ru/ncc/ncc-backend/pkg/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInformings_PrepareMessageList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := zero.NewLogger()

	sms := mocks.NewMockSmsProvider(ctrl)

	informingsRepo := mocks2.NewMockInformings(ctrl)
	informingsRepo.EXPECT().GetEnabled().Times(1).Return(fixtures.FakeInformings(), nil)
	informingsRepo.EXPECT().SetState(gomock.Any(), gomock.Any()).Times(1).Return(nil)
	informingsRepo.EXPECT().SetStart(gomock.Any(), gomock.Any()).Times(2).Return(nil)

	informingLogRepo := mocks2.NewMockInformingLog(ctrl)

	informingsTestCustomersRepo := mocks2.NewMockInformingsTestCustomers(ctrl)
	informingsTestCustomersRepo.EXPECT().Get().Times(1).Return(fixtures.FakeInformingsTestCustomers(), nil)

	customerRepo := mocks2.NewMockCustomers(ctrl)
	customerRepo.EXPECT().Get(gomock.Any()).Times(1).Return(fixtures.FakeCustomers(), nil)

	informings := NewInformings(
		log,
		sms,
		informingsRepo,
		informingLogRepo,
		informingsTestCustomersRepo,
		customerRepo,
	)
	list, err := informings.PrepareMessageList()
	assert.NoError(t, err)
	assert.NotNil(t, list)
	assert.Equal(t, 3, len(list))
}
