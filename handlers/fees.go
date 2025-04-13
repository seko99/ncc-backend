package handlers

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/errors"
	repository2 "code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type FeeLogResponse struct {
	Date          string `json:"date"`
	Amount        string `json:"amount"`
	DepositBefore string `json:"deposit_before"`
}

type Fees struct {
	customersRepo repository2.Customers
	feeRepo       repository2.Fees
}

func (s *Fees) GetFeeLog(c echo.Context) error {
	login, ok := c.Get("login").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, errors.ErrorResponse("f7dad729-c7ee-4345-8daa-23ea8ada2844", "invalid context"))
	}

	customer, err := s.customersRepo.GetByLogin(login)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.ErrorResponse("645cd723-369e-4b2f-879e-926fb87b6cd2"))
	}

	//todo: get period from request
	start, _ := time.Parse("2006-01-02", "2022-02-01")
	fees, err := s.feeRepo.GetByCustomer(customer.Id, repository2.TimePeriod{
		Start: start,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.ErrorResponse("8e3ad85e-9c57-49e3-ae97-2e3e61a2f05b"))
	}

	var result []FeeLogResponse

	for _, f := range fees {
		result = append(result, FeeLogResponse{
			Date:          f.FeeTimestamp.Format("2006-01-02"),
			Amount:        fmt.Sprintf("%0.2f", f.FeeAmount),
			DepositBefore: fmt.Sprintf("%0.2f", f.PrevDeposit),
		})
	}

	return c.JSON(http.StatusOK, result)
}

func NewFees(customersRepo repository2.Customers, feeRepo repository2.Fees) *Fees {
	return &Fees{
		customersRepo: customersRepo,
		feeRepo:       feeRepo,
	}
}
