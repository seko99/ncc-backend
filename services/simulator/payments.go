package simulator

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"fmt"
	"math/rand"
	"time"
)

func (ths *Simulator) createPayments(firstPayment time.Time) error {
	ths.log.Info("Creating payments starting from %v", firstPayment)

	customers, err := ths.customerRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get customers: %w", err)
	}

	paymentTypes, err := ths.paymentTypesRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get payment types: %w", err)
	}

	pid := 100_000

	for idx, c := range customers {
		err := ths.customerRepo.Update(models2.CustomerData{
			CommonData: models2.CommonData{
				Id: c.Id,
			},
			Deposit: 0.0,
		})
		if err != nil {
			ths.log.Error("Can't update customer: %v", err)
		}

		ths.log.Info("Creating payments for customer %d/%d...", idx+1, len(customers))

		currentDate := firstPayment
		for {
			customer, err := ths.customerRepo.GetById(c.Id)
			if err != nil {
				ths.log.Error("Can't get customer: %v", err)
				break
			}

			amount := customer.ServiceInternet.Fee
			paymentTypeId := paymentTypes[rand.Intn(len(paymentTypes))].Id

			err = ths.paymentsRepo.Create(models2.PaymentData{
				Pid:           pid,
				PaymentTypeId: models2.NewNullUUID(paymentTypeId),
				CustomerId:    models2.NewNullUUID(customer.Id),
				Date:          currentDate,
				Amount:        amount,
				DepositBefore: customer.Deposit,
			})
			if err != nil {
				ths.log.Error("Can't create payment: %v", err)
			}

			err = ths.customerRepo.Update(models2.CustomerData{
				CommonData: models2.CommonData{
					Id: customer.Id,
				},
				Deposit: customer.Deposit + amount,
			})
			if err != nil {
				ths.log.Error("Can't update customer: %v", err)
			}

			pid++

			currentDate = currentDate.Add(25 * 24 * time.Hour)
			if currentDate.After(time.Now()) {
				break
			}
		}
	}

	return nil
}
