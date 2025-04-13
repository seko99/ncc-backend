package repository

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

//go:generate mockgen -destination=mocks/mock_sorm_ip_numbering_repo.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/repository SormIpNumbering

type SormIpNumbering interface {
	Create(data models.SormIpNumberingData) error
	Upsert(data models.SormIpNumberingData) error
	Update(data models.SormIpNumberingData) error
	Delete(id string) error
	Get() ([]models.SormIpNumberingData, error)
}
