package handlers

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/errors"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ProfileAddress struct {
	City   string `json:"city"`
	Street string `json:"street"`
	Build  string `json:"build"`
	Flat   string `json:"flat"`
}

type ProfileService struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Fee   string `json:"fee"`
	State string `json:"state"`
}

type ProfileResponse struct {
	Login          string `json:"login"`
	Balance        string `json:"balance"`
	Credit         string `json:"credit"`
	Blocked        bool   `json:"blocked"`
	Phone          string `json:"phone"`
	Scores         int    `json:"scores"`
	Email          string `json:"email"`
	CreditExpire   string `json:"credit_expire"`
	CreditDaysLeft int    `json:"credit_days_left"`

	Address  ProfileAddress   `json:"address"`
	Services []ProfileService `json:"services"`
}

type Profile struct {
	customersRepo repository.Customers
}

func (s *Profile) GetProfile(c echo.Context) error {
	login, ok := c.Get("login").(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, errors.ErrorResponse("ca4dbbf2-554c-4d32-b81b-10dd115a947e", "invalid context"))
	}

	customer, err := s.customersRepo.GetByLogin(login)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errors.ErrorResponse("6c3b979a-8f4c-4830-af16-f673342828e0"))
	}

	return c.JSON(http.StatusOK, &ProfileResponse{
		Login:          customer.Login,
		Balance:        fmt.Sprintf("%0.2f", customer.Deposit),
		Credit:         fmt.Sprintf("%0.2f", customer.Credit),
		Blocked:        s.customersRepo.IsBlocked(customer.BlockingState),
		Phone:          customer.Phone,
		Scores:         customer.Scores,
		Email:          customer.Email,
		CreditExpire:   customer.CreditExpire.Time.Format("2006-01-02"),
		CreditDaysLeft: customer.CreditDaysLeft,
		Address: ProfileAddress{
			City:   customer.City.Name,
			Street: customer.Street.Name,
			Build:  customer.Build,
			Flat:   customer.Flat,
		},
		Services: []ProfileService{
			{
				Type:  models.ServiceTypeInternet,
				Name:  customer.ServiceInternet.Name,
				Fee:   fmt.Sprintf("%0.2f", customer.ServiceInternet.Fee),
				State: models.ServiceState[customer.ServiceInternetState],
			},
			{
				Type:  models.ServiceTypeIptv,
				Name:  customer.ServiceIptv.Name,
				Fee:   fmt.Sprintf("%0.2f", customer.ServiceIptv.Fee),
				State: models.ServiceState[customer.ServiceIptvState],
			},
		},
	})
}

func NewProfile(customersRepo repository.Customers) *Profile {
	return &Profile{
		customersRepo: customersRepo,
	}
}
