package handlers

import (
	"bytes"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	repository2 "code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"net/http"
	"time"
)

type PaymentsByDaysResponse struct {
	Count  int64
	Amount int64
}

type Reports struct {
	log       logger.Logger
	customers repository2.Customers
	payments  repository2.Payments
}

func NewReports(customers repository2.Customers, payments repository2.Payments) Reports {
	return Reports{
		customers: customers,
		payments:  payments,
	}
}

func (r *Reports) Execute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dateStart, err := time.Parse("2006-01-02", ctx.Param("start"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		dateEnd, err := time.Parse("2006-01-02", ctx.Param("end"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		if dateEnd.Before(dateStart) {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "end before start",
			})
			return
		}

		fname := "report.xlsx"
		e := excelize.NewFile()
		idx, err := e.NewSheet("report")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		e.SetActiveSheet(idx)

		row := 1
		dt := dateStart
		for dt.Before(dateEnd.Add(time.Second)) {
			err = e.SetCellStr("report", fmt.Sprintf("A%d", row), fmt.Sprintf("%0.2d.%0.2d.%d", dt.Day(), dt.Month(), dt.Year()))
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			dt = dt.Add(24 * time.Hour)
			row++
		}

		var b bytes.Buffer
		err = e.Write(&b)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Header("Content-Description", "File Transfer")
		ctx.Header("Content-Disposition", "attachment; filename="+fname)
		ctx.Data(http.StatusOK, "application/octet-stream", b.Bytes())
	}
}
