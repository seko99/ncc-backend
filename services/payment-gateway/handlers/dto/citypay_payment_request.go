package dto

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type CitypayPaymentRequestDTO struct {
	Timestamp       time.Time
	IP              string
	QueryType       string
	Amount          float64
	TransactionID   string
	TransactionDate string
	Account         string
	PayElementID    string
	PaymentSystemID string
	Token           string
}

func NewCitypayPaymentRequestDTO() CitypayPaymentRequestDTO {
	return CitypayPaymentRequestDTO{}
}

func (p *CitypayPaymentRequestDTO) Parse(ctx *gin.Context) error {
	var err error

	p.Timestamp = time.Now()

	p.IP = ctx.Query("IP")
	if len(p.IP) == 0 {
		p.IP = ctx.RemoteIP()
	}

	p.QueryType = ctx.Query("QueryType")
	if len(p.QueryType) == 0 {
		return fmt.Errorf("wrong QueryType")
	}

	p.TransactionID = ctx.Query("TransactionId")
	p.TransactionDate = ctx.Query("TransactionDate")
	p.PayElementID = ctx.Query("PayElementId")

	p.PaymentSystemID = ctx.Param("payment_system_id")
	if len(p.PaymentSystemID) == 0 {
		return fmt.Errorf("wrong payment system")
	}

	p.Token = ctx.Param("token")
	if len(p.Token) == 0 {
		return fmt.Errorf("wrong token")
	}

	p.Account = ctx.Query("Account")
	if len(p.Account) == 0 {
		return fmt.Errorf("wrong Account")
	}

	amount := ctx.Query("Amount")

	if len(amount) > 0 {
		p.Amount, err = strconv.ParseFloat(amount, 64)
		if err != nil {
			return fmt.Errorf("wrong Amount")
		}
	}

	return nil
}
