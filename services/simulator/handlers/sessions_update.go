package handlers

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SessionsUpdateEndpoint struct {
	log logger.Logger
	uc  interfaces.SessionsUpdateUsecase
}

func (ths *SessionsUpdateEndpoint) Execute() gin.HandlerFunc {
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

func NewSessionsUpdateEndpoint(log logger.Logger, uc interfaces.SessionsUpdateUsecase) SessionsUpdateEndpoint {
	return SessionsUpdateEndpoint{
		log: log,
		uc:  uc,
	}
}
