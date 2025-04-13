package handlers

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FakeDataClearEndpoint struct {
	log logger.Logger
	uc  interfaces.FakeDataClearUsecase
}

func (ths *FakeDataClearEndpoint) Execute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := dto.NewFakeDataClearEndpointRequest()

		err := ctx.BindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("syntax error: %v", err),
			})
			return
		}

		if err := req.Validate(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("syntax error: %v", err),
			})
			return
		}

		err = ths.uc.Execute(dto.FakeDataClearUsecaseRequest{}.FromFakeDataClearEndpointRequest(req))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"result": "success",
		})
	}
}

func NewFakeDataClearEndpoint(log logger.Logger, uc interfaces.FakeDataClearUsecase) FakeDataClearEndpoint {
	return FakeDataClearEndpoint{
		log: log,
		uc:  uc,
	}
}
