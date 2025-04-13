package handlers

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RadiusStopAllEndpoint struct {
	log logger.Logger
	cfg *config.Config
	uc  interfaces.RadiusStopAllUsecase
}

func (ths *RadiusStopAllEndpoint) Execute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := ths.uc.Execute(dto.RadiusUsecaseRequest{
			Secret:        ths.cfg.Radius.Test.Secret,
			NasIP:         ths.cfg.Radius.Test.Nas.Ip,
			NasIdentifier: ths.cfg.Radius.Test.Nas.Identifier,
			Auth:          ths.cfg.Radius.Test.Auth,
			Acct:          ths.cfg.Radius.Test.Acct,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"result": "success",
		})
	}
}

func NewRadiusStopAllEndpoint(cfg *config.Config, log logger.Logger, uc interfaces.RadiusStopAllUsecase) RadiusStopAllEndpoint {
	return RadiusStopAllEndpoint{
		cfg: cfg,
		log: log,
		uc:  uc,
	}
}
