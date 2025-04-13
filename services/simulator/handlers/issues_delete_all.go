package handlers

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IssuesDeleteAllEndpoint struct {
	log logger.Logger
	uc  interfaces.IssuesDeleteAllUsecase
}

func (ths *IssuesDeleteAllEndpoint) Execute() gin.HandlerFunc {
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

func NewIssuesDeleteAllEndpoint(log logger.Logger, uc interfaces.IssuesDeleteAllUsecase) IssuesDeleteAllEndpoint {
	return IssuesDeleteAllEndpoint{
		log: log,
		uc:  uc,
	}
}
