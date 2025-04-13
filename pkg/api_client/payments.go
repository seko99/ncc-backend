package api_client

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/api_client/dto"
	"fmt"
	"net/http"
)

type PaymentClient struct {
	DefaultClient
}

func (c PaymentClient) Check(request dto.CheckRequestDTO) (*dto.CheckResponseDTO, error) {
	response := &dto.CheckResponseDTO{}
	err := c.Do(http.MethodPost, "payment/check", response, request)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	return response, nil
}

func (c PaymentClient) Payment(request dto.PaymentRequestDTO) (*dto.PaymentResponseDTO, error) {
	response := &dto.PaymentResponseDTO{}
	err := c.Do(http.MethodPost, "payment", response, request)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	return response, nil
}

func (c PaymentClient) GetPaymentSystemByID(id string) (*dto.PaymentSystemResponseDTO, error) {
	response := &dto.PaymentSystemResponseDTO{}
	err := c.Do(http.MethodGet, "payment/system/by/id/"+id, response)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	return response, nil
}

func NewPaymentClient(
	apiURL string,
	apiToken string,
) PaymentClient {
	return PaymentClient{
		DefaultClient: DefaultClient{
			apiURL:   apiURL,
			apiToken: apiToken,
		},
	}
}
