package handlers

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"net/http"
)

type LoginRequest struct {
	Login    string
	Password string
}

type LoginResponse struct {
	Token string
}

type Auth struct {
	log  logger.Logger
	auth repository.Auth
}

func (s *Auth) Login(c *gin.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds LoginRequest

		err := c.Bind(&creds)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "")
		}

		customer, err := s.auth.VerifyLogin(creds.Login, creds.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err)
		}

		token, err := s.auth.CreateToken(customer)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "")
		}

		c.JSON(http.StatusOK, &LoginResponse{
			Token: token,
		})
	}
}

func (s *Auth) Middleware(next echo.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		token, err := s.auth.ExtractToken(authHeader)
		if err != nil {
			s.log.Error(err.Error())
			return
		}

		err = s.auth.VerifyToken(token)
		if err != nil {
			return
		}

		claims, err := s.auth.GetClaims(c)
		if err != nil {
			return
		}

		c.Set("login", claims.Login)

		c.Next()

		return
	}
}

func NewAuth(log logger.Logger, auth repository.Auth) *Auth {
	return &Auth{
		log:  log,
		auth: auth,
	}
}
