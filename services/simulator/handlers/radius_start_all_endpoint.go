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

type RadiusStartAllEndpoint struct {
	log logger.Logger
	cfg *config.Config
	uc  interfaces.RadiusStartAllUsecase
}

func (ths *RadiusStartAllEndpoint) Execute() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		req := dto.RadiusUsecaseRequest{}

		err := ctx.BindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("syntax error: %v", err),
			})
			return
		}

		req.Secret = ths.cfg.Radius.Test.Secret
		req.NasIP = ths.cfg.Radius.Test.Nas.Ip
		req.NasIdentifier = ths.cfg.Radius.Test.Nas.Identifier
		req.Auth = ths.cfg.Radius.Test.Auth
		req.Acct = ths.cfg.Radius.Test.Acct

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
			"response": response,
		})
	}
}

func NewRadiusStartAllEndpoint(cfg *config.Config, log logger.Logger, uc interfaces.RadiusStartAllUsecase) RadiusStartAllEndpoint {
	return RadiusStartAllEndpoint{
		cfg: cfg,
		log: log,
		uc:  uc,
	}
}
