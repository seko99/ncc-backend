package api_client

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/api_client/dto"
	"fmt"
	"net/http"
)

type CustomerClient struct {
	DefaultClient
}

func (c CustomerClient) GetByUID(uid string) (*dto.CustomerResponseDTO, error) {
	customer := &dto.CustomerResponseDTO{}
	err := c.Do(http.MethodGet, "customer/by/uid/"+uid, customer)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	return customer, nil
}

func NewCustomerClient(
	apiURL string,
	apiToken string,
) CustomerClient {
	return CustomerClient{
		DefaultClient: DefaultClient{
			apiURL:   apiURL,
			apiToken: apiToken,
		},
	}
}
