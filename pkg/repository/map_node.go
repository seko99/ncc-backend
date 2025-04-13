package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_map_nodes_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository MapNodes

type MapNodes interface {
	Create(data models.MapNodeData) error
	Update(data models.MapNodeData) error
	Delete(id string) error
	Get() ([]models.MapNodeData, error)
}
