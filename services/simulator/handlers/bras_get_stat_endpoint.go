package handlers

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BrasGetStatEndpoint struct {
	log logger.Logger
	cfg *config.Config
	uc  interfaces.BrasGetStatUsecase
}

func (ths *BrasGetStatEndpoint) Execute() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		req := dto.BrasGetStatUsecaseRequest{}

		response, err := ths.uc.Execute(req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"result": "success",
			"stat":   response,
		})
	}
}

func NewBrasGetStatEndpoint(cfg *config.Config, log logger.Logger, uc interfaces.BrasGetStatUsecase) BrasGetStatEndpoint {
	return BrasGetStatEndpoint{
		cfg: cfg,
		log: log,
		uc:  uc,
	}
}
