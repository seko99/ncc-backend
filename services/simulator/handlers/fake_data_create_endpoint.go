package handlers

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FakeDataCreateEndpoint struct {
	log logger.Logger
	uc  interfaces.FakeDataCreateUsecase
}

func (ths *FakeDataCreateEndpoint) Execute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := dto.NewFakeDataCreateEndpointRequest()

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

		ths.log.Info("request: %+v", req)

		err = ths.uc.Execute(dto.FakeDataCreateUsecaseRequest{}.FromFakeDataCreateEndpointRequest(req))
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

func NewFakeDataCreateEndpoint(log logger.Logger, uc interfaces.FakeDataCreateUsecase) FakeDataCreateEndpoint {
	return FakeDataCreateEndpoint{
		log: log,
		uc:  uc,
	}
}
