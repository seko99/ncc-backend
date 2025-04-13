package handlers

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SessionsDropEndpoint struct {
	log logger.Logger
	uc  interfaces.SessionsDropUsecase
}

func (ths *SessionsDropEndpoint) Execute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := ths.uc.Execute()
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

func NewSessionsDrop(log logger.Logger, uc interfaces.SessionsDropUsecase) SessionsDropEndpoint {
	return SessionsDropEndpoint{
		log: log,
		uc:  uc,
	}
}
