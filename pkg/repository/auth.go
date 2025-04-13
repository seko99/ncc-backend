package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type CustomClaims struct {
	ID        uuid.UUID `json:"id"`
	Login     string    `json:"login"`
	UID       string    `json:"uid"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (p *CustomClaims) Valid() error {
	return nil
}

type Auth interface {
	VerifyLogin(login, password string) (*models.CustomerData, error)
	GetClaims(c *gin.Context) (*CustomClaims, error)
	CreateToken(customer *models.CustomerData) (string, error)
	VerifyToken(token string) error
	ExtractToken(a string) (string, error)
}
