package handlers

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BrasGetSessionsEndpoint struct {
	log logger.Logger
	cfg *config.Config
	uc  interfaces.BrasGetSessionsUsecase
}

func (ths *BrasGetSessionsEndpoint) Execute() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		req := dto.BrasGetSessionsUsecaseRequest{}

		err := ctx.BindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("syntax error: %v", err),
			})
			return
		}

		ths.log.Info("Request: %+v", req)

		response, err := ths.uc.Execute(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"result":   "success",
			"count":    len(response.Sessions),
			"sessions": response.Sessions,
		})
	}
}

func NewBrasGetSessionsEndpoint(cfg *config.Config, log logger.Logger, uc interfaces.BrasGetSessionsUsecase) BrasGetSessionsEndpoint {
	return BrasGetSessionsEndpoint{
		cfg: cfg,
		log: log,
		uc:  uc,
	}
}
