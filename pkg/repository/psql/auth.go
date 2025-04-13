package psql

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Auth struct {
	cfg     *config.Config
	storage *psqlstorage.Storage
}

func (s *Auth) VerifyLogin(login, password string) (*models.CustomerData, error) {
	cr := NewCustomers(s.storage, nil)

	customer, err := cr.GetByLogin(login)
	if err != nil {
		return nil, err
	}

	if customer.Password == password {
		return customer, nil
	}

	return nil, fmt.Errorf("invalid login/password for %s", login)
}

func (s *Auth) GetClaims(c *gin.Context) (*repository.CustomClaims, error) {
	token, err := s.ExtractToken(c.Request.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	jwtToken, err := jwt.ParseWithClaims(token, &repository.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(s.cfg.JwtSecret), nil
	})

	if claims, ok := jwtToken.Claims.(*repository.CustomClaims); ok && claims != nil {
		return claims, nil
	}

	return nil, fmt.Errorf("can't get claims")
}

func (s *Auth) ExtractToken(a string) (string, error) {
	str := strings.Split(a, " ")
	if len(str) == 2 {
		return str[1], nil
	}
	return "", fmt.Errorf("invalid token: %s", a)
}

func (s *Auth) CreateToken(customer *models.CustomerData) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &repository.CustomClaims{
		ID:        uuid.New(),
		Login:     customer.Login,
		UID:       customer.Id,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(48 * time.Hour),
	})

	token, err := t.SignedString([]byte(s.cfg.JwtSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Auth) VerifyToken(token string) error {
	jwtToken, err := jwt.ParseWithClaims(token, &repository.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(s.cfg.JwtSecret), nil
	})

	if err != nil {
		return err
	}

	_ = jwtToken

	return nil
}

func NewAuth(cfg *config.Config, storage *psqlstorage.Storage) *Auth {
	return &Auth{
		cfg:     cfg,
		storage: storage,
	}
}
