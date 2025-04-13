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

type RadiusKillSessionsEndpoint struct {
	log logger.Logger
	cfg *config.Config
	uc  interfaces.RadiusKillSessionsUsecase
}

func (ths *RadiusKillSessionsEndpoint) Execute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := dto.NewRadiusKillSessionsEndpointRequest()

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

		response, err := ths.uc.Execute(dto.RadiusKillSessionsUsecaseRequest{
			Sessions: req.Sessions,
			Random:   req.Random,
		})
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

func NewRadiusKillSessionsEndpoint(cfg *config.Config, log logger.Logger, uc interfaces.RadiusKillSessionsUsecase) RadiusKillSessionsEndpoint {
	return RadiusKillSessionsEndpoint{
		cfg: cfg,
		log: log,
		uc:  uc,
	}
}
