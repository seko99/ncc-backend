package handlers

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InitDictionariesEndpoint struct {
	log logger.Logger
	uc  interfaces.InitDictionariesUsecase
}

func (ths *InitDictionariesEndpoint) Execute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := ths.uc.Execute()
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

func NewInitDictionariesEndpoint(log logger.Logger, uc interfaces.InitDictionariesUsecase) InitDictionariesEndpoint {
	return InitDictionariesEndpoint{
		log: log,
		uc:  uc,
	}
}
