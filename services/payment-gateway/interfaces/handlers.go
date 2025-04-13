package interfaces

import "github.com/gin-gonic/gin"

type PaymentHandler interface {
	Execute() gin.HandlerFunc
}
