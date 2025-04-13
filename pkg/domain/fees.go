package domain

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

type FeeService struct {
	Service interface{}
	Fee     float64
}

type Fee struct {
	FeeLog               models.FeeLogData
	Login                string
	ServiceInternetState int
	NewDeposit           float64
	NewState             int
	NewCreditDaysLeft    int
}
