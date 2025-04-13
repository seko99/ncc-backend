package interfaces

import "code.evixo.ru/ncc/ncc-backend/pkg/api_client/dto"

type CustomerClient interface {
	GetByUID(uid string) (*dto.CustomerResponseDTO, error)
}
