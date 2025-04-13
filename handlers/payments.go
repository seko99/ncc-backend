package handlers

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/errors"
	repository2 "code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type PaymentsResponse struct {
	Date          string `json:"date"`
	Amount        string `json:"amount"`
	DepositBefore string `json:"deposit_before"`
}

type Payments struct {
	customersRepo repository2.Customers
	paymentsRepo  repository2.Payments
}

func (s *Payments) GetPayments(c echo.Context) error {
	login, ok := c.Get("login").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, errors.ErrorResponse("43bb4b19-9802-469c-93f8-2df79904c8fd", "invalid context"))
	}

	customer, err := s.customersRepo.GetByLogin(login)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.ErrorResponse("9bf851c8-2639-4162-8778-b3475ce78352"))
	}

	//todo: get period from request
	start, _ := time.Parse("2006-01-02", "2022-02-01")
	payments, err := s.paymentsRepo.GetPaymentsByCustomer(customer.Id, repository2.TimePeriod{
		Start: start,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.ErrorResponse("b016fcc5-8a59-4112-aa19-6807ea441a8f"))
	}

	var result []PaymentsResponse

	for _, p := range payments {
		result = append(result, PaymentsResponse{
			Date:          p.Date.Format("2006-01-02"),
			Amount:        fmt.Sprintf("%0.2f", p.Amount),
			DepositBefore: fmt.Sprintf("%0.2f", p.DepositBefore),
		})
	}

	return c.JSON(http.StatusOK, result)
}

func NewPayments(customersRepo repository2.Customers, paymentsRepo repository2.Payments) *Payments {
	return &Payments{
		customersRepo: customersRepo,
		paymentsRepo:  paymentsRepo,
	}
}
