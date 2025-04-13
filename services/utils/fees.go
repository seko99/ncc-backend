package utils

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/api_client"
	"code.evixo.ru/ncc/ncc-backend/pkg/api_client/dto"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	repository2 "code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"fmt"
	"time"
)

const (
	PaymentTypeForCompensation = "a1621dcc-fc79-a32a-4947-06e107cd8e91"
)

type FeesUtils struct {
	log           *zero.Logger
	customerRepo  repository2.Customers
	paymentsRepo  repository2.Payments
	feeRepo       repository2.Fees
	paymentClient api_client.PaymentClient
}

func NewFeesUtils(
	cfg *config.Config,
	log *zero.Logger,
	customerRepo repository2.Customers,
	paymentsRepo repository2.Payments,
	feeRepo repository2.Fees,
) *FeesUtils {
	return &FeesUtils{
		log:           log,
		customerRepo:  customerRepo,
		paymentsRepo:  paymentsRepo,
		feeRepo:       feeRepo,
		paymentClient: api_client.NewPaymentClient(cfg.API.URL, cfg.API.Token),
	}
}

type CompensateData struct {
	Customer   models.CustomerData
	Duplicates int
	Amount     float64
}

func (f *FeesUtils) CompensateDuplicateFees(start, end time.Time, dry bool) error {

	fees, err := f.feeRepo.Get(
		repository2.TimePeriod{
			Start: start,
			End:   end,
		},
	)
	if err != nil {
		return fmt.Errorf("can't get fees: %w", err)
	}

	customerFeeMap := map[string][]models.FeeLogData{}

	duplicateFees := []models.FeeLogData{}

	for _, fee := range fees {
		if fee.FeeAmount <= 0.0 {
			continue
		}
		customerID := fee.CustomerId.UUID.String()
		if _, ok := customerFeeMap[customerID]; !ok {
			customerFeeMap[customerID] = []models.FeeLogData{fee}
		} else {
			t := customerFeeMap[customerID][len(customerFeeMap[customerID])-1].FeeTimestamp
			prev := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
			next := prev.Add(24 * time.Hour)
			if fee.FeeTimestamp.After(prev) && fee.FeeTimestamp.Before(next) {
				duplicateFees = append(duplicateFees, fee)
			} else {
				customerFeeMap[customerID] = append(customerFeeMap[customerID], fee)
			}
		}
	}

	customersToCompensate := map[string]CompensateData{}

	for _, dupFee := range duplicateFees {
		f.log.Info("Duplicate fee: login=%s date=%v amount=%0.2f", dupFee.Customer.Login, dupFee.FeeTimestamp, dupFee.FeeAmount)
		customerID := dupFee.CustomerId.UUID.String()
		if _, ok := customersToCompensate[customerID]; !ok {
			customersToCompensate[customerID] = CompensateData{
				Customer:   dupFee.Customer,
				Duplicates: 1,
				Amount:     dupFee.FeeAmount,
			}
		} else {
			data := customersToCompensate[customerID]
			customersToCompensate[customerID] = CompensateData{
				Customer:   data.Customer,
				Duplicates: data.Duplicates + 1,
				Amount:     data.Amount + dupFee.FeeAmount,
			}
		}
	}

	totalAmount := 0.0
	for _, c := range customersToCompensate {

		if !dry {
			f.log.Info("Compensate payment: %s (%d) amount=%0.2f", c.Customer.Login, c.Duplicates, c.Amount)
			err := f.makePayment(c.Customer, c.Amount)
			if err != nil {
				f.log.Error("Can't make payment: %v", err)
			}
		} else {
			f.log.Info("[DRY] Compensate payment: %s (%d) amount=%0.2f", c.Customer.Login, c.Duplicates, c.Amount)
		}

		totalAmount += c.Amount
	}
	f.log.Info("Summary: %d payments amount=%0.2f", len(customersToCompensate), totalAmount)

	return nil
}

func (f *FeesUtils) makePayment(customer models.CustomerData, amount float64) error {
	_, err := f.paymentClient.Payment(dto.PaymentRequestDTO{
		CreatedBy:      "compensate-duplicate-fees",
		CustomerID:     customer.Id,
		PaymentTypeID:  PaymentTypeForCompensation,
		Amount:         amount,
		ChargeDailyFee: false,
		ClearCredit:    false,
		Descr:          fmt.Sprintf("корректировка ошибочно начисленной абонентской платы"),
	})
	return err
}
