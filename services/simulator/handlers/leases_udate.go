package handlers

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LeasesUpdateEndpoint struct {
	log logger.Logger
	uc  interfaces.LeasesUpdateUsecase
}

func (ths *LeasesUpdateEndpoint) Execute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := ths.uc.Execute()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
		}

		ctx.JSON(http.StatusOK, map[string]interface{}{
			"result": "success",
		})
	}
}

func NewLeasesUpdateEndpoint(log logger.Logger, uc interfaces.LeasesUpdateUsecase) LeasesUpdateEndpoint {
	return LeasesUpdateEndpoint{
		log: log,
		uc:  uc,
	}
}
