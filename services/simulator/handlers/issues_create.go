package handlers

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IssuesCreateEndpoint struct {
	log logger.Logger
	uc  interfaces.IssuesCreateUsecase
}

func (ths *IssuesCreateEndpoint) Execute() gin.HandlerFunc {
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

func NewIssuesCreateEndpoint(log logger.Logger, uc interfaces.IssuesCreateUsecase) IssuesCreateEndpoint {
	return IssuesCreateEndpoint{
		log: log,
		uc:  uc,
	}
}
