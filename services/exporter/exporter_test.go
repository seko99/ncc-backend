package exporter

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/fixtures"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	mocks2 "code.evixo.ru/ncc/ncc-backend/pkg/repository/mocks"
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	mocks3 "code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExporter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := zero.NewLogger()

	writer := mocks3.NewMockExportWriter(ctrl)
	writer.EXPECT().Write(gomock.Any()).Times(1).Return(nil)
	writer.EXPECT().GetErrors(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return([]exporter.ExportData{}, nil)

	customerRepo := mocks2.NewMockCustomers(ctrl)
	customerRepo.EXPECT().Get(gomock.Any()).Times(1).Return(fixtures.FakeCustomers(), nil)

	paymentRepo := mocks2.NewMockPayments(ctrl)
	paymentTypesRepo := mocks2.NewMockPaymentTypes(ctrl)
	serviceInternetRepo := mocks2.NewMockServiceInternet(ctrl)
	documentTypesRepo := mocks2.NewMockDocumentTypes(ctrl)
	ipNumberingRepo := mocks2.NewMockSormIpNumbering(ctrl)
	gatewayRepo := mocks2.NewMockSormGateway(ctrl)

	sormCustomersRepo := mocks2.NewMockSormCustomers(ctrl)
	sormCustomersRepo.EXPECT().Get().Times(1).Return([]models.SormCustomersData{
		{
			IP:    "172.18.0.10",
			Login: "test",
		},
	}, nil)
	sormCustomersRepo.EXPECT().Upsert(gomock.Any()).Times(1).Return(nil)

	sormCustomersErrorsRepo := mocks2.NewMockSormCustomersErrors(ctrl)
	sormCustomersErrorsRepo.EXPECT().DeleteAll().Times(1).Return(nil)

	sormCustomerServicesRepo := mocks2.NewMockSormCustomerServices(ctrl)

	sormCustomerServicesErrorsRepo := mocks2.NewMockSormCustomerServicesErrors(ctrl)

	sormExportStatusRepo := mocks2.NewMockSormExportStatus(ctrl)
	sormExportStatusRepo.EXPECT().Upsert(gomock.Any()).Times(1).Return(nil)

	exporterService := NewExporter(
		log,
		writer,
		customerRepo,
		paymentRepo,
		paymentTypesRepo,
		serviceInternetRepo,
		documentTypesRepo,
		ipNumberingRepo,
		gatewayRepo,
		sormCustomersRepo,
		sormCustomersErrorsRepo,
		sormCustomerServicesRepo,
		sormCustomerServicesErrorsRepo,
		sormExportStatusRepo,
	)

	customerList, err := exporterService.exportCustomers()
	assert.NoError(t, err)
	assert.Equal(t, len(fixtures.FakeCustomers()), len(customerList))
}
